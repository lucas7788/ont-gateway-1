package io

import (
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/define"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// BuyerBuyDtokenInput ...
type BuyerBuyDtokenInput struct {
	Tx string
}

// BuyerBuyDtokenOutput ...
type BuyerBuyDtokenOutput struct {
	io2.BaseResp
	EndpointTokens []EndpointToken
}

type EndpointToken struct {
	Token    Token
	Endpoint string
}

type Token struct {
	TokenTemplate define.TokenTemplate
	Buyer         common.Address
	OnchainItemId string
}
