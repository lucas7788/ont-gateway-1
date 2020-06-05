package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ddxf"
	common2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
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
const (
	sellerCollectionName   = "seller"
	endpointCollectionName = "sellerendpoint"
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

// Init for this collection
func initDb() (err error) {
	opts := &options.IndexOptions{}
	opts.SetName("u-seller")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "seller", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	_, err = instance.MongoOfficial().Collection(sellerCollectionName).Indexes().CreateOne(context.Background(), index)
	return
}

func SaveDataMetaService(input io.SellerSaveDataMetaInput, ontId string) (output io.SellerSaveDataMetaOutput) {
	// verify hash.
	h, err := ddxf.HashObject(input.DataMeta)
	if err != nil || hex.EncodeToString(h[:]) != input.DataMetaHash {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	dataMetaHash, err := common.Uint256FromHexString(input.DataMetaHash)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	dataHash, err := common.Uint256FromHexString(input.DataMetaHash)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	dataId := common2.GenerateUUId(config.UUID_PRE_DATAID)
	info := DataIdInfo{
		DataId:       dataId,
		DataType:     input.ResourceType,
		DataMetaHash: dataMetaHash,
		DataHash:     dataHash,
		Owner:        ontId,
	}
	txHash, err := instance.OntSdk().DefaultDataIdContract().Invoke(ServerAccount, "registerDataId", []interface{}{info.ToBytes()})
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	event, err := instance.OntSdk().GetSmartCodeEvent(txHash.ToHexString())
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	if event.State != 1 {
		output.Code = http.StatusInternalServerError
		output.Msg = fmt.Sprintf("registerDataId failed, txHash: %s", txHash.ToHexString())
		return
	}
	// reg identity.
	dataStore := &io.SellerSaveDataMeta{
		DataMeta:     input.DataMeta,
		DataMetaHash: input.DataMetaHash,
		ResourceType: input.ResourceType,
		OntId:        ontId,
		DataIds:      dataId,
		Fee:          input.Fee,
		Stock:        input.Stock,
		ExpiredDate:  input.ExpiredDate,
		DataEndpoint: input.DataEndpoint,
	}

	// store meta hash id.
	err = InsertElt(DataMetaCollection, dataStore)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}

	return
}

func SaveTokenMetaService(input io.SellerSaveTokenMetaInput, ontId string) (output io.SellerSaveTokenMetaOutput) {
	// verify hash.
	h, err := ddxf.HashObject(input.TokenMeta)

	if err != nil || hex.EncodeToString(h[:]) != input.TokenMetaHash {
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
	//TODO send mp
	_, _, data, err := forward.JSONRequest("POST", input.MPEndpoint, mpParamBs)
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
	output.Result = "SUCCESS"
	fmt.Println(ets)
	return
}
