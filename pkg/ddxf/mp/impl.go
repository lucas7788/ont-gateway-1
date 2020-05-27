package mp

import (
	"context"
	"encoding/hex"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-go-sdk"
	signature2 "github.com/ontio/ontology/core/signature"
	"github.com/ontio/ontology/core/types"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
	"time"
)

const (
	mpCollectionName       = "marketplace"
	endpointCollectionName = "mpendpoint"
)

type MarketplaceImpl struct {
	endpointImpl Endpoint
}

func NewMarketplaceImpl(acc *ontology_go_sdk.Account) *MarketplaceImpl {
	return &MarketplaceImpl{
		endpointImpl: NewEndpointImpl(acc),
	}
}

type EndpointImpl struct {
	mpAccount *ontology_go_sdk.Account
}

func NewEndpointImpl(acc *ontology_go_sdk.Account) *EndpointImpl {
	return &EndpointImpl{
		mpAccount: acc,
	}
}

func (this *MarketplaceImpl) AddRegistry(input io.MPAddRegistryInput) (output io.MPAddRegistryOutput) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := hex.DecodeString(input.PubKey)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	_, err = instance.MongoOfficial().Collection(mpCollectionName).InsertOne(ctx, input)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func (this *MarketplaceImpl) RemoveRegistry(input io.MPRemoveRegistryInput) (output io.MPRemoveRegistryOutput) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ad := &io.MPAddRegistryInput{}
	filter := bson.M{"mp": input.MP}
	err := instance.MongoOfficial().Collection(mpCollectionName).FindOne(ctx, filter).Decode(ad)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	pkBs, err := hex.DecodeString(ad.PubKey)
	if err != nil {
		panic("")
	}
	pk, err := keypair.DeserializePublicKey(pkBs)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	err = signature2.Verify(pk, []byte(input.MP), input.Sign)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	filter = bson.M{"mp": input.MP}
	_, err = instance.MongoOfficial().Collection(mpCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func (this *MarketplaceImpl) Endpoint() Endpoint {
	return this.endpointImpl
}

// Init for this collection
func (m *MarketplaceImpl) Init() (err error) {

	opts := &options.IndexOptions{}
	opts.SetName("u-mp")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "mp", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	_, err = instance.MongoOfficial().Collection(mpCollectionName).Indexes().CreateOne(context.Background(), index)
	return
}

func (this *EndpointImpl) GetAuditRule(io.MPEndpointGetAuditRuleInput) (output io.MPEndpointGetAuditRuleOutput) {
	return
}

func (this *EndpointImpl) GetFee(io.MPEndpointGetFeeInput) (output io.MPEndpointGetFeeOutput) {
	output.Count = 10
	output.Type = io.ONG
	return
}

func (this *EndpointImpl) GetChallengePeriod(io.MPEndpointGetChallengePeriodInput) (output io.MPEndpointGetChallengePeriodOutput) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ou := &io.MPEndpointGetChallengePeriodOutput{}
	err := instance.MongoOfficial().Collection(endpointCollectionName).FindOne(ctx, nil).Decode(ou)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func (this *EndpointImpl) GetItemMetaSchema(io.MPEndpointGetItemMetaSchemaInput) (output io.MPEndpointGetItemMetaSchemaOutput) {

	return
}

func (this *EndpointImpl) GetItemMeta(io.MPEndpointGetItemMetaInput) (output io.MPEndpointGetItemMetaOutput) {
	instance.MongoOfficial().Collection(endpointCollectionName).FindOne(context.Background(), nil)
	return
}

func (this *EndpointImpl) QueryItemMetas(io.MPEndpointQueryItemMetasInput) (output io.MPEndpointQueryItemMetasOutput) {
	opts := options.Find()
	// Sort by `_id` field descending
	opts.SetSort(bson.D{bson.E{Key: "_id", Value: -1}})
	cursor, err := instance.MongoOfficial().Collection(endpointCollectionName).Find(context.Background(), nil, opts)
	if err != nil {
		return
	}

	itemMetas := make([]map[string]interface{}, 0)
	err = cursor.All(context.Background(), itemMetas)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	output.ItemMetas = itemMetas
	return
}

func (this *EndpointImpl) PublishItemMeta(input io.MPEndpointPublishItemMetaInput) (output io.MPEndpointPublishItemMetaOutput) {
	txBs, err := hex.DecodeString(input.SignedDDXFTx)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	tx, err := types.TransactionFromRawBytes(txBs)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	muTx, err := tx.IntoMutable()
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = instance.MongoOfficial().Collection(endpointCollectionName).InsertOne(ctx, input)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	err = instance.OntSdk().GetKit().SignToTransaction(muTx, this.mpAccount)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	txHash, err := instance.OntSdk().GetKit().SendTransaction(muTx)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	state, err := getSmartCodeEvent(txHash.ToHexString())
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	if state == 0 {
		output.Code = http.StatusInternalServerError
		output.Msg = "tx failed"
	}
	return
}

func getSmartCodeEvent(txHash string) (byte, error) {
	instance.OntSdk().GetKit().WaitForGenerateBlock(30 * time.Second)
	event, err := instance.OntSdk().GetKit().GetSmartContractEvent(txHash)
	if err != nil {
		return 0, err
	}
	return event.State, nil
}
