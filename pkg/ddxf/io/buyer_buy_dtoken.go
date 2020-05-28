package io

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/define"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// BuyerBuyDtokenInput ...
type BuyerBuyDtokenInput struct {
	Tx            string
	OnchainItemID string //resource_id
}

// BuyerBuyDtokenOutput ...
type BuyerBuyDtokenOutput struct {
	io2.BaseResp
	EndpointTokens EndpointTokens
}
type EndpointTokens struct {
	Tokens   []define.TokenTemplate `bson:"tokens" json:"tokens"`
	Endpoint string                 `bson:"endpoint" json:"endpoint"`
}
