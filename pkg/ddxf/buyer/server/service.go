package server

import (
	"context"
	"encoding/json"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
)

const (
	buyerCollectionName = "buyer"
	buyDToken           = "buyDToken"
	useTokenM           = "useToken"
)

func Init() error {
	opts := &options.IndexOptions{}
	opts.SetName("u-qrCodeId")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "qrCodeId", Value: bsonx.Int32(1)}},
		Options: opts,
	}
	_, err := instance.MongoOfficial().Collection(buyerCollectionName).Indexes().CreateOne(context.Background(), index)
	return err
}

func insertOne(data interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := instance.MongoOfficial().Collection(buyerCollectionName).InsertOne(ctx, data)
	return err
}
func insertMany(data []interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := instance.MongoOfficial().Collection(buyerCollectionName).InsertMany(ctx, data)
	return err
}

func findOne(filter bson.M, data interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return instance.MongoOfficial().Collection(buyerCollectionName).FindOne(ctx, filter).Decode(data)
}

func BuyDTokenService(param io.BuyerBuyDtokenInput) (output io.BuyerBuyDtokenOutput) {
	txHash, err := sendTx(param.SignedTx)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	output.EndpointTokens, err = HandleEvent(txHash, "buyDtoken")
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	p := make([]interface{}, len(output.EndpointTokens))
	for i := 0; i < len(output.EndpointTokens); i++ {
		p[i] = output.EndpointTokens[i]
	}
	err = insertMany(p)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}

func UseTokenService(input io.BuyerUseTokenInput) (output io.BuyerUseTokenOutput) {
	txHash, err := sendTx(input.Tx)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	endpointTokens, err := HandleEvent(txHash, "useToken")
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	paramBs, err := json.Marshal(endpointTokens)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	//向seller发请求
	_, _, data, err := forward.JSONRequest("useToken", "", paramBs)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	output.Result = data
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = instance.MongoOfficial().Collection(buyerCollectionName).InsertOne(ctx, output.Result)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}
