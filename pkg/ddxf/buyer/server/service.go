package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
	"strings"
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

func findOne(filter bson.M, data interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return instance.MongoOfficial().Collection(buyerCollectionName).FindOne(ctx, filter).Decode(data)
}

func BuyDTokenService(param io.BuyerBuyDtokenInput) (output io.BuyerBuyDtokenOutput) {
}

func UseTokenService(input io.BuyerUseTokenInput) (output io.BuyerUseTokenOutput) {
	code, err := qrCode.BuildBuyQrCode("testnet", input.OnchainItemId, input.N, input.Buyer)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	err = insertOne(code)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	return qrCode.BuildBuyGetQrCodeRsp(code.QrCodeId), nil

	endpointTokens, err := sendTxAndGetTokens(input.Tx, "useToken")
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
