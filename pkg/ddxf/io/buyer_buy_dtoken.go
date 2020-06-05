package io

import (
	"github.com/ontio/ontology/common"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
)

// BuyerBuyDtokenInput ...
type BuyerBuyDtokenInput struct {
	SignedTx string
}

// BuyerBuyDtokenOutput ...
type BuyerBuyDtokenOutput struct {
	io2.BaseResp
	EndpointTokens []EndpointToken
}

type EndpointToken struct {
	Token    Token  `bson:"token" json:"token"`
	Endpoint string `bson:"endpoint" json:"endpoint"`
}

type Token struct {
	TokenTemplate param.TokenTemplate  `bson:"token_template" json:"token_template"`
	Buyer         common.Address `bson:"buyer" json:"buyer"`
	OnchainItemId string         `bson:"onchain_item_id" json:"onchain_item_id"`
}
