package server

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"crypto/sha256"

	"fmt"

	"strings"

	"github.com/kataras/go-errors"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology-go-sdk"
	common3 "github.com/ontio/ontology-go-sdk/common"
	common2 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/ontid"
	"github.com/zhiqiangxu/ddxf"
	config2 "github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/middleware/jwt"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateOntIdService(input GenerateOntIdInput) (output GenerateOntIdOutput) {
	output.ReqID = input.ReqID
	input.UserId = input.Party + input.UserId
	// callback(output)
	// return
	var err error
	defer func() {
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
			fmt.Println("err", err)
		}
		callback(output)
	}()

	ui := UserInfo{}
	filter := bson.M{"user_id": input.UserId}
	err = FindElt(UserInfoCollection, filter, &ui)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println(err)
		return
	}
	if err == nil {
		output.OntId = ui.OntId
		return
	}
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	account := GetAccount(input.UserId)
	ontid := "did:ont:" + account.Address.ToBase58()
	tx, err := instance.DDXFSdk().GetOntologySdk().Native.OntId.NewRegIDWithPublicKeyTransaction(config2.GasPrice,
		config2.GasLimit, ontid, account.PublicKey)
	if err != nil {
		return
	}
	err = instance.DDXFSdk().SignTx(tx, account)
	if err != nil {
		return
	}
	txHash, err := common.SendRawTx(tx)
	if err != nil && strings.Contains(err.Error(), "already registered") {
		ui.OntId = ontid
		ui.UserId = input.UserId
		err = InsertElt(UserInfoCollection, ui)
		fmt.Println(err)
		return
	}
	if err != nil {
		txHas := tx.Hash()
		fmt.Printf("userId: %s,ontid: %s, txHash: %s\n", input.UserId, ontid, txHas.ToHexString())
		return
	}
	var evt *common3.SmartContactEvent
	evt, err = instance.DDXFSdk().GetSmartCodeEvent(txHash)
	if err != nil {
		return
	}
	if evt.State != 1 {
		err = errors.New("event state is not 1")
		return
	}

	ui.OntId = ontid
	ui.UserId = input.UserId
	err = InsertElt(UserInfoCollection, ui)
	if err == nil {
		output.OntId = ontid
	}

	return
}

func PublishService(input PublishInput) (output PublishOutput) {
	input.UserID = input.Party + input.UserID
	output.ReqID = input.ReqID
	var err error
	defer func() {
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
			fmt.Println("error: ", err)
		}
		callback(output)
	}()

	// 抽取openKGID
	openKGID := input.OpenKGID
	filter := bson.M{"openkg_id": openKGID}
	param := PublishInput{}
	var findError error
	findError = FindElt(OpenKgPublishParamCollection, filter, &param)
	if findError != nil && findError != mongo.ErrNoDocuments {
		err = findError
		return
	}
	seller := GetAccount(input.UserID)
	ontID := "did:ont:" + seller.Address.ToBase58()
	jwtToken, err := jwt.GenerateJwt(ontID)
	if err != nil {
		return
	}

	headers := map[string]string{
		"Authorization": jwtToken,
	}

	if input.Delete {
		// 处理删除
		var (
			tx       *types.MutableTransaction
			iMutTx   *types.Transaction
			bs, data []byte
		)
		tx, err = instance.DDXFSdk().DefMpKit().BuildDeleteTx([]byte(param.OnChainId))
		if err != nil {
			return
		}
		err = instance.DDXFSdk().SignTx(tx, seller)
		if err != nil {
			return
		}

		iMutTx, err = tx.IntoImmutable()
		if err != nil {
			return
		}

		param := server.DeleteParam{SignedTx: hex.EncodeToString(common2.SerializeToBytes(iMutTx))}
		bs, err = json.Marshal(param)
		if err != nil {
			return
		}
		_, _, data, err = forward.PostJSONRequestWithRetry(config.SellerUrl+server.FreezeUrl, bs, headers, 10)
		if err != nil {
			return
		}

		res := server.DeleteOutput{}
		err = json.Unmarshal(data, &res)
		if err != nil {
			return
		}
		err = res.Error()
		return
	}

	// 1. 抽取data meta
	dataMetas := input.Datas
	dataMetaHashArray := make([]string, len(dataMetas))
	for i := 0; i < len(dataMetas); i++ {
		if dataMetas[i]["url"] == nil {
			err = errors.New("url empty")
			return
		}
		var hash [sha256.Size]byte
		hash, err = ddxf.HashObject(dataMetas[i])
		if err != nil {
			return
		}
		dataMetaHashArray[i] = string(hash[:])
	}
	getDataIdParam := server.GetDataIdParam{
		DataMetaHashArray: dataMetaHashArray,
	}
	var data, paramBs []byte
	paramBs, err = json.Marshal(getDataIdParam)
	if err != nil {
		return
	}
	//查询哪些data id需要上链
	_, _, data, err = forward.PostJSONRequestWithRetry(config.SellerUrl+server.GetDataIdByDataMetaHashUrl, paramBs, headers, 10)
	if err != nil {
		return
	}

	res := make(map[string]interface{})
	if data != nil {
		err = json.Unmarshal(data, &res)
		if err != nil {
			return
		}
	}
	ones := make([]io.DataMetaOne, 0)
	for i := 0; i < len(dataMetas); i++ {
		var dataMetaHash [sha256.Size]byte
		dataMetaHash, err = ddxf.HashObject(dataMetas[i])
		if err != nil {
			err = fmt.Errorf("1 HashObject error: %s", err)
			return
		}
		var hash common2.Uint256
		hash, err = common2.Uint256ParseFromBytes(dataMetaHash[:])
		if err != nil {
			err = fmt.Errorf("2 Uint256ParseFromBytes error: %s", err)
			return
		}
		dataId := res[hash.ToHexString()]
		if dataId == nil {
			dataId := common.GenerateOntId()
			g := &ontid.Group{
				Members:   []interface{}{},
				Threshold: 1,
			}
			signers := []ontid.Signer{ontid.Signer{
				Id:    []byte(dataId),
				Index: 1,
			}}
			tx, err := instance.DDXFSdk().GetOntologySdk().Native.OntId.NewRegIDWithControllerTransaction(config2.GasPrice,
				config2.GasLimit, dataId, g, signers)
			if err != nil {
				return
			}
			err = instance.DDXFSdk().SignTx(tx, seller)
			if err != nil {
				return
			}
			imutTx, err := tx.IntoImmutable()
			if err != nil {
				return
			}
			//send tx to seller
			ri := server.RegisterOntIdInput{
				SignedTx: hex.EncodeToString(common2.SerializeToBytes(imutTx)),
			}
			data, err := json.Marshal(ri)
			if err != nil {
				return
			}
			code, _, _, err := forward.PostJSONRequest(config.SellerUrl+server.RegisterOntId, data, nil)
			if err != nil {
				return
			}
			txHash := tx.Hash()
			if code != http.StatusOK {
				err = fmt.Errorf("register ontid tx failed, txHash: %s", txHash.ToHexString())
				return
			}

			one := io.DataMetaOne{
				DataMeta:     dataMetas[i],
				DataMetaHash: hash.ToHexString(),
				DataEndpoint: config.SellerUrl,
				ResourceType: 0,
				DataHash:     common2.UINT256_EMPTY.ToHexString(),
				DataId:       dataId,
			}
			ones = append(ones, one)
			res[hash.ToHexString()] = dataId
		}
	}
	// invoke seller saveDataMeta
	for i := 0; i < len(ones); i++ {
		var (
			txMut  *types.MutableTransaction
			iMutTx *types.Transaction
		)
		var hash, dataHash common2.Uint256
		hash, err = common2.Uint256FromHexString(ones[i].DataMetaHash)
		if err != nil {
			err = fmt.Errorf("3 Uint256FromHexString error: %s", err)
			return
		}
		dataHash, err = common2.Uint256FromHexString(ones[i].DataHash)
		if err != nil {
			err = fmt.Errorf("3 Uint256FromHexString error: %s", err)
			return
		}
		attr := &ontology_go_sdk.DDOAttribute{
			Key:       []byte("DataMetaHash"),
			Value:     hash[:],
			ValueType: []byte{},
		}
		attr2 := &ontology_go_sdk.DDOAttribute{
			Key:       []byte("DataHash"),
			Value:     dataHash[:],
			ValueType: []byte{},
		}
		attrs := make([]*ontology_go_sdk.DDOAttribute, 0)
		attrs = append(attrs, attr)
		attrs = append(attrs, attr2)

		txMut, err = instance.OntSdk().GetKit().Native.OntId.NewAddAttributesTransaction(
			500, 2000000, ontID, attrs, seller.PublicKey,
		)
		if err != nil {
			return
		}
		txMut.Payer = payer.Address
		err = instance.DDXFSdk().SignTx(txMut, payer)
		if err != nil {
			return
		}
		err = instance.DDXFSdk().SignTx(txMut, seller)
		if err != nil {
			return
		}
		iMutTx, err = txMut.IntoImmutable()
		if err != nil {
			return
		}
		txHash := txMut.Hash()
		fmt.Println("txhash:", txHash.ToHexString())
		ones[i].SignedTx = hex.EncodeToString(common2.SerializeToBytes(iMutTx))
	}
	saveDataMetaArray := io.SellerSaveDataMetaArrayInput{
		DataMetaOneArray: ones,
	}
	fmt.Printf("ontID: %s, userId: %s\n", ontID, input.UserID)

	var bs []byte
	bs, err = json.Marshal(saveDataMetaArray)
	if err != nil {
		return
	}
	start := time.Now().Unix()
	_, _, data, err = forward.PostJSONRequestWithRetry(config2.SellerUrl+server.SaveDataMetaArrayUrl, bs, headers, 10)
	end := time.Now().Unix()
	fmt.Printf("openkg send seller SaveDataMetaArrayUrl cost time: %d\n", end-start)
	if err != nil {
		return
	}

	templates := make([]*market_place_contract.TokenTemplate, 0)
	for i := 0; i < len(dataMetas); i++ {
		var dataMetaHash [sha256.Size]byte
		dataMetaHash, err = ddxf.HashObject(dataMetas[i])
		if err != nil {
			return
		}
		u, _ := common2.Uint256ParseFromBytes(dataMetaHash[:])
		dataId := res[u.ToHexString()]
		tt := &market_place_contract.TokenTemplate{
			DataID:     dataId.(string),
			TokenHashs: []string{"1"},
			Endpoint:   "aaaa",
		}
		templates = append(templates, tt)
	}

	fmt.Println("******templates:", templates)

	resourceId := common.GenerateUUId(config.UUID_RESOURCE_ID)
	input.OnChainId = resourceId
	// 2. save data metas and Publish item
	// send request to seller
	var itemMetaHash [32]byte
	itemMetaHash, err = ddxf.HashObject(input.Item)
	ddo := market_place_contract.ResourceDDO{
		Manager:      seller.Address,
		ItemMetaHash: itemMetaHash,
	}

	addr, _ := common2.AddressFromHexString(config2.OEP4ContractAddr)
	item := market_place_contract.DTokenItem{
		Fee: market_place_contract.Fee{
			ContractAddr: addr,
			ContractType: split_policy_contract.OEP4,
			Count:        1,
		},
		ExpiredDate: uint64(time.Now().Unix()) + uint64(time.Hour*24*30),
		Stocks:      10000,
		Templates:   templates,
	}
	split := split_policy_contract.SplitPolicyRegisterParam{
		AddrAmts: []*split_policy_contract.AddrAmt{
			&split_policy_contract.AddrAmt{
				To:      seller.Address,
				Percent: 100,
			},
		},
		TokenTy: split_policy_contract.ONG,
	}
	var (
		tx *types.MutableTransaction
	)
	if findError == mongo.ErrNoDocuments {
		tx, err = instance.DDXFSdk().DefMpKit().BuildPublishTx([]byte(resourceId), ddo, item, split)
		if err != nil {
			return
		}
		err = instance.DDXFSdk().SignTx(tx, seller)
		if err != nil {
			return
		}
	} else {
		tx, err = instance.DDXFSdk().DefMpKit().BuildUpdateTx([]byte(resourceId), ddo, item, split)
		if err != nil {
			return
		}
		err = instance.DDXFSdk().SignTx(tx, seller)
		if err != nil {
			return
		}
	}

	iMutTx, err := tx.IntoImmutable()
	if err != nil {
		return
	}

	publishInput := io.MPEndpointPublishItemMetaInput{
		SignedDDXFTx: hex.EncodeToString(common2.SerializeToBytes(iMutTx)),
		ItemMeta:     io.PublishItemMeta{ItemMeta: input.Item, OnchainItemID: resourceId},
		MPEndpoint:   config2.MPEndpoint,
	}
	bs, err = json.Marshal(publishInput)
	if err != nil {
		return
	}

	// send req to seller
	start = time.Now().Unix()
	_, _, data, err = forward.PostJSONRequestWithRetry(config2.SellerUrl+server.PublishMPItemMetaUrl, bs, headers, 10)
	end = time.Now().Unix()
	fmt.Printf("openkg publish service cost time: %d\n", end-start)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	var publishOutput io.SellerPublishMPItemMetaOutput
	err = json.Unmarshal(data, &publishOutput)
	if err != nil {
		return
	}

	err = publishOutput.Error()
	if err != nil {
		return
	}
	err = InsertElt(OpenKgPublishParamCollection, input)

	return
}

func deleteService(input DeleteInput) (output DeleteOutput) {
	var err error
	defer func() {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		output.ReqID = input.ReqID
	}()
	input.UserID = input.Party + input.UserID
	user := GetAccount(input.UserID)
	tx, err := instance.DDXFSdk().DefMpKit().BuildDeleteTx([]byte(input.ResourceID))
	if err != nil {
		return
	}
	instance.DDXFSdk().SignTx(tx, user)
	imut, err := tx.IntoImmutable()
	if err != nil {
		return
	}
	di := server.DeleteInput{
		SignedTx: hex.EncodeToString(common2.SerializeToBytes(imut)),
	}
	data, err := json.Marshal(di)
	if err != nil {
		return
	}
	//send tx to seller
	forward.PostJSONRequest(config.SellerUrl+server.DeleteUrl, data, nil)
	return
}

func regDataService(input RegDataInput) (output AddAttributesOutput) {
	GetAccount(input.PartyDataID)
	return
}

func deleteAttributesService(input DeleteAttributesInput) (output AddAttributesOutput) {
	var err error
	defer func() {
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
		}
	}()
	input.UserID = input.Party + input.UserID
	acc := GetAccount(input.UserID)
	ontId := config2.PreOntId + acc.Address.ToBase58()
	key := []byte("DataMetaHash")
	err = deleteAttribute(ontId, key, acc)
	if err != nil {
		return
	}
	key2 := []byte("DataHash")
	err = deleteAttribute(ontId, key2, acc)
	return
}

func deleteAttribute(ontId string, key []byte, acc *ontology_go_sdk.Account) (err error) {
	tx, err := instance.DDXFSdk().GetOntologySdk().Native.OntId.NewRemoveAttributeTransaction(config2.GasPrice,
		config2.GasLimit, ontId, key, acc.PublicKey)
	if err != nil {
		return
	}
	err = instance.DDXFSdk().SignTx(tx, payer)
	if err != nil {
		return
	}
	err = instance.DDXFSdk().SignTx(tx, acc)
	if err != nil {
		return
	}
	txHash, err := common.SendRawTx(tx)
	if err != nil {
		return
	}
	evt, err := instance.DDXFSdk().GetSmartCodeEvent(txHash)
	if err != nil {
		return
	}
	if evt.State != 1 {
		err = fmt.Errorf("tx failed, txhash: %s", txHash)
		return
	}
	return
}

func buyAndUseService(input BuyAndUseInput) (output BuyAndUseOutput) {
	// callback(output)
	// return
	input.UserID = input.Party + input.UserID
	output.ReqID = input.ReqID
	defer func() {
		callback(output)
	}()

	user := GetAccount(input.UserID)

	ontID := "did:ont:" + user.Address.ToBase58()
	jwtToken, err := jwt.GenerateJwt(ontID)
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}

	headers := map[string]string{
		"Authorization": jwtToken,
	}

	filter := bson.M{"openkg_id": input.OpenKGID}
	param := PublishInput{}
	err = FindElt(OpenKgPublishParamCollection, filter, &param)
	if param.Delete {
		output.Msg = "has delete"
		output.Code = http.StatusInternalServerError
		return
	}
	tokenTemplate := market_place_contract.TokenTemplate{
		DataID:     input.DataID,
		TokenHashs: []string{},
	}
	var tx *types.MutableTransaction
	contractAddr, _ := common2.AddressFromHexString(config2.BuyAndUseContractAddr)
	contract := instance.DDXFSdk().DefContract(contractAddr)
	tx, err = contract.BuildTx("buyAndUseToken",
		[]interface{}{[]byte(param.OnChainId), 1, user.Address, payer.Address, tokenTemplate.ToBytes()})
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}
	err = instance.DDXFSdk().SignTx(tx, user)
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}
	err = instance.DDXFSdk().SignTx(tx, payer)
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}

	iMutTx, err := tx.IntoImmutable()
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}
	bs, err := json.Marshal(io.BuyerBuyAndUseDtokenInput{SignedTx: hex.EncodeToString(common2.SerializeToBytes(iMutTx))})
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}

	// send req to seller
	_, _, data, err := forward.PostJSONRequest(config.SellerUrl+server.BuyAndUseDTokenUrl, bs, headers)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	var buyAndUseOutput io.BuyerBuyAndUseDtokenOutput
	err = json.Unmarshal(data, &buyAndUseOutput)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	err = buyAndUseOutput.Error()
	return
}
