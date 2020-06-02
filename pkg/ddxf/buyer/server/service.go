package server

import (
	"context"
	"encoding/json"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"net/http"
)

const (
	buyDToken           = "buyDToken"
	useToken            = "useToken"
)

func BuyDTokenService(param io.BuyerBuyDtokenInput) (output io.BuyerBuyDtokenOutput) {
	txHash, err := instance.OntSdk().SendTx(param.SignedTx)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	instance.OntSdk().WaitForGenerateBlock()
	output.EndpointTokens, err = HandleEvent(txHash, buyDToken)
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

	txHash, err := instance.OntSdk().SendTx(input.Tx)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	instance.OntSdk().WaitForGenerateBlock()
	endpointTokens, err := HandleEvent(txHash, useToken)
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
	_, _, data, err := forward.JSONRequest("useToken", input.TokenOpEndpoint, paramBs)
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
