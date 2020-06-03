package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"net/http"
)

const (
	buyDToken = "buyDToken"
	useToken  = "useToken"
)

func BuyDTokenService(param io.BuyerBuyDtokenInput) (output io.BuyerBuyDtokenOutput) {
	txHash, err := instance.OntSdk().SendTx(param.SignedTx)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	output.EndpointTokens, err = common.HandleEvent(txHash, buyDToken)
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
	pk := keypair.SerializePublicKey(BuyerMgrAccount.PublicKey)
	sellerUseTokenParam := io.SellerTokenLookupEndpointUseTokenInput{
		Tx:             input.Tx,
		BuyerOntId:     "",
		BuyerPublicKey: hex.EncodeToString(pk),
	}
	paramBs, err := json.Marshal(sellerUseTokenParam)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	fmt.Println(string(paramBs))
	//向seller发请求
	_, _, data, err := forward.PostJSONRequest(input.TokenOpEndpoint, paramBs)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	fmt.Println("888888888888888")
	output.Result = data

	type SellerReturn struct {
		Result interface{} `bson:"result",json:"result"`
	}
	rr := SellerReturn{
		Result: string(data),
	}
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = instance.MongoOfficial().Collection(buyerCollectionName).InsertOne(ctx, rr)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}
