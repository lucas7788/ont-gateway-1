package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/core/types"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/registry/client"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
	"time"
)

const (
	endpointCollectionName = "mp_endpoint"
)

var MpAccount *ontology_go_sdk.Account

func AddRegistryService(input io.MPAddRegistryInput) (output io.MPAddRegistryOutput) {
	output = io.MPAddRegistryOutput(client.Sdk().AddEndpoint(io.RegistryAddEndpointInput(input)))
	return
}

func RemoveRegistryService(input io.MPRemoveRegistryInput) (output io.MPRemoveRegistryOutput) {
	output = io.MPRemoveRegistryOutput(client.Sdk().RemoveEndpoint(io.RegistryRemoveEndpointInput(input)))
	return
}

// Init for this collection
func Init() (err error) {
	opts := &options.IndexOptions{}
	opts.SetName("u-item-meta_id")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "item-meta_id", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	_, err = instance.MongoOfficial().Collection(endpointCollectionName).Indexes().CreateOne(context.Background(), index)
	return
}

func GetAuditRuleService(io.MPEndpointGetAuditRuleInput) (output io.MPEndpointGetAuditRuleOutput) {
	return
}

func GetFeeService(io.MPEndpointGetFeeInput) (output io.MPEndpointGetFeeOutput) {
	output.Count = 10
	output.Type = io.ONG
	return
}

func GetChallengePeriodService(io.MPEndpointGetChallengePeriodInput) (output io.MPEndpointGetChallengePeriodOutput) {
	output.Period = 7 * 24 * 3600 * time.Second
	return
}

func GetItemMetaSchemaService(io.MPEndpointGetItemMetaSchemaInput) (output io.MPEndpointGetItemMetaSchemaOutput) {
	output.ItemMetaSchema = map[string]interface{}{
		"@context": map[string]interface{}{
			"sec":        "http://purl.org/security#",
			"xsd":        "http://www.w3.org/2001/XMLSchema#",
			"rdf":        "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
			"dc":         "http://purl.org/dc/terms/",
			"sec:signer": map[string]interface{}{"@type": "@id"},
			"dc:created": map[string]interface{}{"@type": "xsd:dateTime"},
		},
		"@id":                "http://example.org/sig1",
		"@type":              []interface{}{"rdf:Graph", "sec:SignedGraph"},
		"dc:created":         "2011-09-23T20:21:34Z",
		"sec:signer":         "http://payswarm.example.com/i/john/keys/5",
		"sec:signatureValue": "doc1",
		"@graph": map[string]interface{}{
			"@id":      "http://example.org/fact1",
			"dc:title": "Hello World!",
		},
	}
	return
}

func GetItemMetaService(input io.MPEndpointGetItemMetaInput) (output io.MPEndpointGetItemMetaOutput) {
	id, err := primitive.ObjectIDFromHex(input.ItemMetaID)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	filter := bson.M{"_id": id}
	item := io.PublishItemMeta{}
	err = instance.MongoOfficial().Collection(endpointCollectionName).FindOne(context.Background(), filter).Decode(&item)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	output.ItemMeta = item
	return
}

func QueryItemMetasService(input io.MPEndpointQueryItemMetasInput) (output io.MPEndpointQueryItemMetasOutput) {
	pageNum := input.PageNum
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize := input.PageSize
	if pageSize < 0 {
		pageSize = 0
	}
	if pageSize > 100 {
		pageSize = 100
	}
	skip := (pageNum - 1) * pageSize
	opts := options.Find()
	opts.Limit = &pageSize
	opts.Skip = &skip
	filter := bson.D{}
	cursor, err := instance.MongoOfficial().Collection(endpointCollectionName).Find(context.Background(), filter, opts)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	itemMetas := make([]io.ItemMeta, 0)
	err = cursor.All(context.Background(), &itemMetas)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	output.ItemMetas = itemMetas
	return
}

func PublishItemMetaService(input io.MPEndpointPublishItemMetaInput) (output io.MPEndpointPublishItemMetaOutput) {
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
	_, err = instance.MongoOfficial().Collection(endpointCollectionName).InsertOne(ctx, input.ItemMeta)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	err = instance.OntSdk().GetKit().SignToTransaction(muTx, MpAccount)
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
	fmt.Println("mp send tx:", txHash.ToHexString())
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
