package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	common2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
)

const (
	buyDToken = "buyDToken"
	useToken  = "useToken"
)

//BuyDTokenService buy token service, POST
//input parameter:
//type BuyerBuyDtokenInput struct {
//	SignedTx string
//}
//this method will send transaction to ontology blockchain,
func BuyDTokenService(param io.BuyerBuyDtokenInput) (output io.BuyerBuyDtokenOutput) {
	txHash, err := common2.SendTx(param.SignedTx)
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
	bt := BuyerToken{
		TxHash: txHash,
		Tokens: output.EndpointTokens,
	}
	err = insertOne(buyerDtokenCollectionName, bt)
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
	//向seller发请求
	fmt.Println("seller param: ", string(paramBs))

	_, _, data, err := forward.PostJSONRequest(input.TokenOpEndpoint+server.UseTokenUrl, paramBs, nil)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	output.Result = data

	type SellerReturn struct {
		Result interface{} `bson:"result",json:"result"`
	}
	rr := SellerReturn{
		Result: string(data),
	}
	err = insertOne(buyerResultCollectionName, rr)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}
