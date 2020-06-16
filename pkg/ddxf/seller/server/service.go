package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ontio/ontology-crypto/signature"
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ddxf"
	common2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/mp/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
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
		return err
	}
	pri, _ := hex.DecodeString("c19f16785b8f3543bbaf5e1dbb5d398dfa6c85aaad54fc9d71203ce83e505c07")
	ServerAccount, _ = ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	Wallet = ontology_go_sdk.NewWallet("./wallet.dat")
	Pwd = []byte("111111")
	return nil
}

func SaveDataMetaService(input io.SellerSaveDataMetaInput, ontId string) (output io.SellerSaveDataMetaOutput) {
	if input.DataMeta["ISDN"] == "" {
		output.Code = http.StatusBadRequest
		output.Msg = "datameta does not contain ISDN"
		return
	}
	// verify hash.
	txHash, err := instance.OntSdk().SendTx(input.SignedTx)
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
		Fee:          input.Fee,
		Stock:        input.Stock,
		ExpiredDate:  input.ExpiredDate,
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

func PublishMPItemMetaService(input io.MPEndpointPublishItemMetaInput, ontId string) (output io.SellerPublishMPItemMetaOutput) {
	mpParamBs, err := json.Marshal(input)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	err = InsertElt(ItemMetaCollectionDdxf, input)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	//TODO send mp
	if input.MPEndpoint != "" {
		_, _, data, err := forward.JSONRequest("POST", input.MPEndpoint+server.PublishItemMeta, mpParamBs, nil)
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
		txHash, err := instance.OntSdk().SendTx(input.SignedDDXFTx)
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
			return
		}
		instance.OntSdk().WaitForGenerateBlock()
		evt, err := instance.OntSdk().GetSmartCodeEvent(txHash)
		if err != nil {
			output.Code = http.StatusInternalServerError
			output.Msg = err.Error()
			return
		}
		if evt.State == 0 {
			output.Code = http.StatusInternalServerError
			output.Msg = "tx failed"
			return
		}
	}
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
	txHash, err := instance.OntSdk().SendRawTx(mutTx)
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

//书籍
func GetDataByOnchainIdService(onchainItemId string, buyer common.Address, template param.TokenTemplate) (interface{}, error) {
	//根据 onchainId 拿到真实的数据
	data := io.SellerSaveDataMeta{}
	find := bson.M{"dataId": template.DataID}
	err := FindElt(DataMetaCollection, find, &data)
	if err != nil {
		return nil, err
	}
	return "http://127.0.0.1:20335/ddxf/storage/download/" + data.DataMeta["downloadParam"].(string), nil
}
