package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ontio/ontology-crypto/signature"
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ddxf"
	common2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/mp/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	WalletName    string
	ServerAccount *ontology_go_sdk.Account
	Wallet        *ontology_go_sdk.Wallet
	Pwd           []byte
	Version       string
)

const (
	ONTAuthScanProtocol = "http://172.29.36.101" + config.SellerPort + "/ddxf/seller/getQrCodeDataByQrCodeId"
	QrCodeCallback      = "http://172.29.36.101" + config.SellerPort + "/ddxf/seller/qrCodeCallbackSendTx"
)

func InitSellerImpl() error {
	err := initDb()
	if err != nil {
		panic(fmt.Sprintf("InitSellerImpl initDb:%v", err))
		return err
	}
	pri, _ := hex.DecodeString("c19f16785b8f3543bbaf5e1dbb5d398dfa6c85aaad54fc9d71203ce83e505c07")
	ServerAccount, _ = ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	Wallet = ontology_go_sdk.NewWallet("./wallet.dat")
	Pwd = []byte("111111")
	return nil
}

func GetDataIdByDataMetaHashService(param GetDataIdParam) (map[string]string, error) {
	res := make(map[string]string)
	for _, item := range param.DataMetaHashArray {
		dataStore := &io.SellerSaveDataMeta{}
		filter := bson.M{"dataMetaHash": item}
		err := FindElt(DataMetaCollection, filter, dataStore)
		if err != nil && err != mongo.ErrNoDocuments {
			fmt.Println("seller: ", err)
			return nil, err
		}
		if err == mongo.ErrNoDocuments {
			continue
		}
		res[dataStore.DataMetaHash] = dataStore.DataId
	}
	return res, nil
}

func SaveDataMetaArrayService(input io.SellerSaveDataMetaArrayInput,
	ontid string) (output io.SellerSaveDataMetaOutput) {
	txHash, err := common2.SendTx(input.SignedTx)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	event, err := instance.OntSdk().GetSmartCodeEvent(txHash)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	if event.State != 1 {
		output.Code = http.StatusInternalServerError
		output.Msg = fmt.Sprintf("registerDataId failed, txHash: %s", txHash)
		return
	}
	for _, item := range input.DataMetaOneArray {
		dataStore := &io.SellerSaveDataMeta{
			DataMeta:     item.DataMeta,
			DataMetaHash: item.DataMetaHash,
			ResourceType: item.ResourceType,
			OntId:        ontid,
			DataId:       item.DataId,
			DataEndpoint: item.DataEndpoint,
		}
		// store meta hash id.
		fmt.Println("dataStore:", dataStore)
		err = InsertElt(DataMetaCollection, dataStore)
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
		}
	}
	return
}

func SaveDataMetaService(input io.SellerSaveDataMetaInput, ontId string) (output io.SellerSaveDataMetaOutput) {
	if input.DataMeta["ISDN"] == "" {
		output.Code = http.StatusBadRequest
		output.Msg = "datameta does not contain ISDN"
		return
	}
	// verify hash.
	txHash, err := common2.SendTx(input.SignedTx)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	fmt.Printf("[seller] saveDataMeta txhash: %s\n", txHash)
	event, err := instance.OntSdk().GetSmartCodeEvent(txHash)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	if event.State != 1 {
		output.Code = http.StatusInternalServerError
		output.Msg = fmt.Sprintf("registerDataId failed, txHash: %s", txHash)
		return
	}
	// reg identity.
	dataStore := &io.SellerSaveDataMeta{
		DataMeta:     input.DataMeta,
		DataMetaHash: input.DataMetaHash,
		ResourceType: input.ResourceType,
		OntId:        ontId,
		DataId:       input.DataId,
		DataEndpoint: input.DataEndpoint,
	}

	fmt.Println("dataStore: ", dataStore)
	// store meta hash id.
	err = InsertElt(DataMetaCollection, dataStore)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	output.DataId = input.DataId
	return
}

func SaveTokenMetaService(input io.SellerSaveTokenMetaInput, ontId string) (output io.SellerSaveTokenMetaOutput) {
	// verify hash.
	h, err := ddxf.HashObject(input.TokenMeta)

	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	if hex.EncodeToString(h[:]) != input.TokenMetaHash {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	adD := &io.SellerSaveDataMeta{}

	filterD := bson.M{"dataMetaHash": input.DataMetaHash, "ontId": ontId}
	err = FindElt(DataMetaCollection, filterD, adD)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	tokenStore := &io.SellerSaveTokenMeta{
		TokenMeta:     input.TokenMeta,
		TokenMetaHash: input.TokenMetaHash,
		DataMetaHash:  input.DataMetaHash,
		TokenEndpoint: input.TokenEndpoint,
		OntId:         ontId,
	}

	err = InsertElt(TokenMetaCollection, tokenStore)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func FreezeService(param DeleteParam, ontId string) (res FreezeOutput) {
	tx, err := utils.TransactionFromHexString(param.SignedTx)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = err.Error()
		return
	}
	mutTx, err := tx.IntoMutable()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = err.Error()
		return
	}
	txHash, err := instance.OntSdk().GetKit().SendTransaction(mutTx)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = err.Error()
		return
	}
	evt, err := instance.DDXFSdk().GetSmartCodeEvent(txHash.ToHexString())
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = err.Error()
		return
	}
	if evt.State != 1 {
		res.Code = http.StatusInternalServerError
		res.Msg = "event state is not 1, txHash: " + txHash.ToHexString()
		return
	}
	return
}

func PublishMPItemMetaService(input io.MPEndpointPublishItemMetaInput, ontId string) (output io.SellerPublishMPItemMetaOutput) {
	mpParamBs, err := json.Marshal(input)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	//TODO send mp
	if input.MPEndpoint != "" {
		start := time.Now().Unix()
		_, _, data, err := forward.PostJSONRequestWithRetry(input.MPEndpoint+server.PublishItemMeta, mpParamBs, nil, 10)
		end := time.Now().Unix()
		fmt.Printf("seller PublishMPItemMetaService cost time: %d\n", end-start)
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
			return
		}
		res := io.MPEndpointPublishItemMetaOutput{}
		err = json.Unmarshal(data, &res)
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
			return
		}
	} else {
		txHash, err := common2.SendTx(input.SignedDDXFTx)
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
			return
		}
		evt, err := instance.DDXFSdk().GetSmartCodeEvent(txHash)
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
			return
		}
		if evt.State != 1 {
			output.Code = http.StatusInternalServerError
			output.Msg = "tx failed"
			return
		}
	}
	err = InsertElt(ItemMetaCollectionDdxf, input)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}

func DeleteService(input DeleteInput) (output DeleteOutput) {
	//TODO
	//send to mp
	bs, err := json.Marshal(input)
	if err != nil {
		return
	}
	_, _, data, err := forward.PostJSONRequest(config.MpUrl+server.Delete, bs, nil)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	output.Result = data
	return
}

func UpdateService(input UpdateInput) (output DeleteOutput) {
	//TODO
	//send to mp
	bs, err := json.Marshal(input)
	if err != nil {
		return
	}
	_, _, data, err := forward.PostJSONRequest(config.MpUrl+server.Update, bs, nil)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	output.Result = data
	return
}

func LookupService(io.SellerDataLookupEndpointLookupInput) (output io.SellerDataLookupEndpointLookupOutput) {
	return
}

func UseTokenService(input io.SellerTokenLookupEndpointUseTokenInput) (output io.SellerTokenLookupEndpointUseTokenOutput) {
	tx, err := utils.TransactionFromHexString(input.Tx)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	mutTx, err := tx.IntoMutable()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	txHash, err := common2.SendRawTx(mutTx)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	ets, err := common2.HandleEvent(txHash, "useToken")
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	//TODO
	if len(ets) == 0 {
		output.Code = http.StatusInternalServerError
		output.Msg = ""
		return
	}
	output.Result, err = GetDataByOnchainIdService(ets[0].Token.OnchainItemId, ets[0].Token.Buyer, ets[0].Token.TokenTemplate)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	fmt.Println(ets)
	return
}

func BuyAndUseDTokenService(input io.BuyerBuyAndUseDtokenInput) (output io.BuyerBuyAndUseDtokenOutput) {
	tx, err := utils.TransactionFromHexString(input.SignedTx)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	mutTx, err := tx.IntoMutable()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	txHash, err := common2.SendRawTx(mutTx)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	buy, use, err := common2.HanleBuyAndUseToken(txHash)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	//TODO
	if len(buy) == 0 || len(use) == 0 {
		output.Code = http.StatusInternalServerError
		output.Msg = ""
		return
	}
	output.EndpointTokens = buy
	output.Result, err = GetDataByOnchainIdService(use[0].Token.OnchainItemId, use[0].Token.Buyer, use[0].Token.TokenTemplate)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	fmt.Println(buy, use)
	return
}

//书籍
func GetDataByOnchainIdService(onchainItemId string, buyer common.Address, template market_place_contract.TokenTemplate) (interface{}, error) {
	//根据 onchainId 拿到真实的数据
	data := io.SellerSaveDataMeta{}
	find := bson.M{"dataId": template.DataID}
	err := FindElt(DataMetaCollection, find, &data)
	if err != nil {
		return nil, err
	}

	if data.DataMeta["url"] != nil {
		return data.DataMeta["url"].(string), nil
	}

	return "http://127.0.0.1:20335/ddxf/storage/download/" + data.DataMeta["downloadParam"].(string), nil
}
