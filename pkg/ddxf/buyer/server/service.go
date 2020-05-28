package server

import (
	"context"
	"errors"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/define"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/http_utils"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
)

const (
	buyerCollectionName = "buyer"
)


func Init() error {
	opts := &options.IndexOptions{}
	opts.SetName("u-tx")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "tx", Value: bsonx.Int32(1)}},
		Options: opts,
	}
	_, err := instance.MongoOfficial().Collection(buyerCollectionName).Indexes().CreateOne(context.Background(), index)
	return err
}

func BuyDtokenService(input io.BuyerBuyDtokenInput) (output io.BuyerBuyDtokenOutput) {
	var err error
	output.EndpointTokens, err = sendTxAndGetTemplates(input.Tx, input.OnchainItemID)
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = instance.MongoOfficial().Collection(buyerCollectionName).InsertOne(ctx, output.EndpointTokens)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}

func UseTokenService(input io.BuyerUseTokenInput) (output io.BuyerUseTokenOutput) {
	endpointTokens, err := sendTxAndGetTemplates(input.Tx, input.OnchainItemID)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	data, err := http_utils.NewHttpClient("").PostRequest(endpointTokens, input.TokenOpEndpoint)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	output.Result = data
	return
}

func sendTxAndGetTemplates(txHex string, OnchainItemID string) (io.EndpointTokens, error) {
	tx, err := utils.TransactionFromHexString(txHex)
	if err != nil {
		return io.EndpointTokens{}, err
	}
	mutTx, err := tx.IntoMutable()
	if err != nil {
		return io.EndpointTokens{}, err
	}
	txHash, err := instance.OntSdk().GetKit().SendTransaction(mutTx)
	if err != nil {
		return io.EndpointTokens{}, err
	}
	event, err := instance.OntSdk().GetSmartCodeEvent(txHash.ToHexString())
	if err != nil {
		return io.EndpointTokens{}, err
	}
	if event.State != 1 {
		return io.EndpointTokens{}, errors.New("tx failed")
	}
	res, err := instance.OntSdk().DDXFContract(0, 0,
		nil).PreInvoke("getTokenTemplates", []interface{}{OnchainItemID})
	if err != nil {
		return io.EndpointTokens{}, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return io.EndpointTokens{}, err
	}
	source := common.NewZeroCopySource(data)
	templates, err := define.DeserializeTokenTemplates(source)
	if err != nil {
		return io.EndpointTokens{}, err
	}
	endpoint, err := define.ReadString(source)
	if err != nil {
		return io.EndpointTokens{}, err
	}
	return io.EndpointTokens{
		Tokens:   templates,
		Endpoint: endpoint,
	}, nil
}
