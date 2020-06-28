package server

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"crypto/sha256"

	"fmt"

	"github.com/kataras/go-errors"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology-go-sdk"
	common3 "github.com/ontio/ontology-go-sdk/common"
	common2 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
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
	"strings"
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
	tx, err := instance.DDXFSdk().GetOntologySdk().Native.OntId.NewRegIDWithPublicKeyTransaction(defGasPrice,
		defGasLimit, ontid, account.PublicKey)
	if err != nil {
		return
	}
	_, err = instance.DDXFSdk().SignTx(tx, account)
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
		tx, err = instance.DDXFSdk().SignTx(tx, seller)
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
		_, _, data, err = forward.PostJSONRequest(config.SellerUrl+server.FreezeUrl, bs, headers)
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
	_, _, data, err = forward.PostJSONRequest(config.SellerUrl+server.GetDataIdByDataMetaHashUrl, paramBs, headers)
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
			dataId := common.GenerateUUId(config.UUID_PRE_DATAID)
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
	attrs := make([]*ontology_go_sdk.DDOAttribute, 0)
	for i := 0; i < len(ones); i++ {
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
		attrs = append(attrs, attr)
		attrs = append(attrs, attr2)
	}
	var (
		txMut  *types.MutableTransaction
		iMutTx *types.Transaction
	)
	fmt.Printf("ontID: %s, userId: %s\n", ontID, input.UserID)
	txMut, err = instance.OntSdk().GetKit().Native.OntId.NewAddAttributesByIndexTransaction(
		500, 2000000, ontID, attrs, 1,
	)
	if err != nil {
		return
	}
	txMut, err = instance.DDXFSdk().SignTx(txMut, seller)
	if err != nil {
		return
	}
	iMutTx, err = txMut.IntoImmutable()
	if err != nil {
		return
	}
	saveDataMetaArray := io.SellerSaveDataMetaArrayInput{
		DataMetaOneArray: ones,
		SignedTx:         hex.EncodeToString(common2.SerializeToBytes(iMutTx)),
	}
	fmt.Println("txhash:", hex.EncodeToString(common2.SerializeToBytes(iMutTx)))

	var bs []byte
	bs, err = json.Marshal(saveDataMetaArray)
	if err != nil {
		return
	}
	_, _, data, err = forward.PostJSONRequest(config2.SellerUrl+server.SaveDataMetaArrayUrl, bs, headers)
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
	InsertElt(OpenKgPublishParamCollection, input)
	// 2. save data metas and Publish item
	// send request to seller
	var itemMetaHash [32]byte
	itemMetaHash, err = ddxf.HashObject(input.Item)
	ddo := market_place_contract.ResourceDDO{
		Manager:      seller.Address,
		ItemMetaHash: itemMetaHash,
	}

	addr, _ := common2.AddressFromHexString("195d72da6725e8243a52803f6de4cd93df48fc1f")

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
		tx, err = instance.DDXFSdk().SignTx(tx, seller)
		if err != nil {
			return
		}
	} else {
		tx, err = instance.DDXFSdk().DefMpKit().BuildUpdateTx([]byte(resourceId), ddo, item, split)
		if err != nil {
			return
		}
		tx, err = instance.DDXFSdk().SignTx(tx, seller)
		if err != nil {
			return
		}
	}

	iMutTx, err = tx.IntoImmutable()
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
	start := time.Now().Unix()
	_, _, data, err = forward.PostJSONRequest(config2.SellerUrl+server.PublishMPItemMetaUrl, bs, headers)
	end := time.Now().Unix()
	fmt.Printf("openkg publish service cost time:%d\n", end-start)
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

	return
}

func deleteService(input DeleteInput) (output DeleteOutput) {
	input.UserID = input.Party + input.UserID
	user := GetAccount(input.UserID)
	tx, err := instance.DDXFSdk().DefMpKit().BuildDeleteTx([]byte(input.ResourceID))
	if err != nil {
		return
	}
	instance.DDXFSdk().SignTx(tx, user)
	//send tx to seller
	forward.PostJSONRequest(config.SellerUrl+server.DeleteUrl, []byte{}, nil)
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
	tx, err = contract.BuildTx("buyAndUseToken", nil,
		[]interface{}{[]byte(param.OnChainId), 1, user.Address, payer.Address, tokenTemplate.ToBytes()})
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}
	tx, err = instance.DDXFSdk().SignTx(tx, user)
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}
	tx, err = instance.DDXFSdk().SignTx(tx, payer)
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
