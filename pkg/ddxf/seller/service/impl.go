package service

import (
	"context"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/contract"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sellerconfig"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
	"strings"
)

const (
	sellerCollectionName   = "marketplace"
	endpointCollectionName = "sellerendpoint"
)

var (
	DefSellerImpl *SellerImpl
)

type SellerImpl struct {
	wallet              *sdk.Wallet
	walletpswd          string
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

	identity, err := self.wallet.NewDefaultSettingIdentity([]byte(self.walletpswd))
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
	}

	// store meta hash id.
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = instance.MongoOfficial().Collection(sellerCollectionName).InsertOne(ctx, dataStore)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}

	return
}

func (self *SellerImpl) SaveTokenMeta(input io.SellerSaveTokenMetaInput, ontId string) (output io.SellerSaveTokenMetaOutput) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// verify hash.
	h, err := ddxf.HashObject(input.TokenMeta)
	if err != nil || h != input.TokenMetaHash {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	adD := &io.SellerSaveDataMeta{}
	filterD := bson.M{"dataMetaHash": input.DataMetaHash, "ontId": ontId}
	err = instance.MongoOfficial().Collection(sellerCollectionName).FindOne(ctx, filterD).Decode(adD)
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

	_, err = instance.MongoOfficial().Collection(sellerCollectionName).InsertOne(ctx, tokenStore)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func (self *SellerImpl) PublishMPItemMeta(input io.SellerPublishMPItemMetaInput, ontId string) (output io.SellerPublishMPItemMetaOutput, response *qrCode.QrCodeResponse) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	adT := &io.SellerSaveTokenMeta{}
	filterT := bson.M{"tokenMetaHash": input.TokenMetaHash, "ontId": ontId}
	err := instance.MongoOfficial().Collection(sellerCollectionName).FindOne(ctx, filterT).Decode(adT)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	adD := &io.SellerSaveDataMeta{}
	filterD := bson.M{"dataMetaHash": input.DataMetaHash, "ontId": ontId}
	err = instance.MongoOfficial().Collection(sellerCollectionName).FindOne(ctx, filterD).Decode(adD)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	arr := strings.Split(ontId, ":")
	if len(arr) != 3 {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	sellerAddress, err := common.AddressFromBase58(arr[2])

	// dataMeta related in data contract tx.
	tokenTemplate := param.TokenTemplate{
		DataIDs:   adD.DataIds,
		TokenHash: adT.TokenMetaHash,
	}
	itemMetaHash, err := ddxf.HashObject(input.ItemMeta)

	resourceIdBytes, rosourceDDOBytes, itemBytes := contract.ConstructPublishParam(sellerAddress, tokenTemplate, adT.TokenEndpoint, itemMetaHash, adD.ResourceType, adD.Fee, adD.ExpiredDate, adD.Stock, adD.DataIds)
	qrCodex, err := qrCode.BuildPublishQrCode(sellerconfig.DefSellerConfig.NetType, input.MPContractHash, resourceIdBytes, rosourceDDOBytes, itemBytes, arr[2], ontId)
	response = qrCode.BuildQrCodeResponse(qrCodex.QrCodeId)
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
}

func (self *DataLookupEndpointImpl) Lookup(io.SellerDataLookupEndpointLookupInput) (output io.SellerTokenLookupEndpointLookupOutput) {
}

type TokenLookupEndpointImpl struct {
}

func (self *TokenLookupEndpointImpl) Lookup(io.SellerTokenLookupEndpointLookupInput) (output io.SellerTokenLookupEndpointLookupOutput) {
}

type TokenOpEndpointImpl struct {
}

func (self *TokenOpEndpointImpl) Lookup(io.SellerTokenLookupEndpointUseTokenInput) (output io.SellerTokenLookupEndpointUseTokenOutput) {
}
