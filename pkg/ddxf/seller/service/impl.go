package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhiqiangxu/ddxf"
	common2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sellerconfig"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sql"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
)

const (
	sellerCollectionName   = "seller"
	endpointCollectionName = "sellerendpoint"
)

var (
	DefSellerImpl *SellerImpl
)

type SellerImpl struct {
	dataLookupEndpoint  DataLookupEndpoint
	tokenLookupEndpoint TokenLookupEndpoint
	tokenOpEndpoint     TokenOpEndpoint
}

func InitSellerImpl() *SellerImpl {
	s := &SellerImpl{
		dataLookupEndpoint:  DataLookupEndpointImpl{},
		tokenLookupEndpoint: TokenLookupEndpointImpl{},
		tokenOpEndpoint:     TokenOpEndpointImpl{},
	}
	s.Init()
	DefSellerImpl = s
	return DefSellerImpl
}

// Init for this collection
func (m *SellerImpl) Init() (err error) {

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

func (self *SellerImpl) SaveDataMeta(input io.SellerSaveDataMetaInput, ontId string) (output io.SellerSaveDataMetaOutput) {
	// verify hash.
	h, err := ddxf.HashObject(input.DataMeta)
	if err != nil || h != input.DataMetaHash {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	identity, err := sellerconfig.DefSellerConfig.Wallet.NewDefaultSettingIdentity([]byte(sellerconfig.DefSellerConfig.Pswd))
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	// reg identity.
	dataStore := &io.SellerSaveDataMeta{
		DataMeta:     input.DataMeta,
		DataMetaHash: input.DataMetaHash,
		ResourceType: input.ResourceType,
		OntId:        ontId,
		DataIds:      identity.ID,
		Fee:          input.Fee,
		Stock:        input.Stock,
		ExpiredDate:  input.ExpiredDate,
		DataEndpoint: input.DataEndpoint,
	}

	// store meta hash id.
	err = sql.InsertElt(sql.DataMetaCollection, dataStore)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}

	return
}

func (self *SellerImpl) SaveTokenMeta(input io.SellerSaveTokenMetaInput, ontId string) (output io.SellerSaveTokenMetaOutput) {
	// verify hash.
	h, err := ddxf.HashObject(input.TokenMeta)
	if err != nil || h != input.TokenMetaHash {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	adD := &io.SellerSaveDataMeta{}

	filterD := bson.M{"dataMetaHash": input.DataMetaHash, "ontId": ontId}
	err = sql.FindElt(sql.DataMetaCollection, filterD, adD)
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

	err = sql.InsertElt(sql.TokenMetaCollection, tokenStore)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func (self *SellerImpl) PublishMPItemMeta(input io.MPEndpointPublishItemMetaInput, ontId string) (output io.SellerPublishMPItemMetaOutput) {
	mpParamBs, err := json.Marshal(input)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	//TODO send mp
	_, _, data, err := forward.JSONRequest("POST", config.PublishItemMetaUrl, mpParamBs)
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

func (self *SellerImpl) DataLookupEndpoint() (output DataLookupEndpoint) {
	return self.dataLookupEndpoint
}

func (self *SellerImpl) TokenLookupEndpoint() (output TokenLookupEndpoint) {
	return self.tokenLookupEndpoint
}

func (self *SellerImpl) TokenOpEndpoint() (output TokenOpEndpoint) {
	return self.tokenOpEndpoint
}

type DataLookupEndpointImpl struct {
	SellerImpl *SellerImpl
}

func (self DataLookupEndpointImpl) Lookup(io.SellerDataLookupEndpointLookupInput) (output io.SellerDataLookupEndpointLookupOutput) {
	return
}

type TokenLookupEndpointImpl struct {
}

func (self TokenLookupEndpointImpl) Lookup(io.SellerTokenLookupEndpointLookupInput) (output io.SellerTokenLookupEndpointLookupOutput) {
	return
}

type TokenOpEndpointImpl struct {
}

func (self TokenOpEndpointImpl) UseToken(input io.SellerTokenLookupEndpointUseTokenInput) (output io.SellerTokenLookupEndpointUseTokenOutput) {
	txHash, err := instance.OntSdk().SendTx(input.Tx)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	instance.OntSdk().WaitForGenerateBlock()
	ets, err := common2.HandleEvent(txHash, "useToken")
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	fmt.Println(ets)
	return
}
