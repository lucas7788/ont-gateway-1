package io

import (
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// BuyerBuyAndUseDtokenInput ...
type BuyerBuyAndUseDtokenInput struct {
	SignedTx string
}

// BuyerBuyAndUseDtokenOutput ...
type BuyerBuyAndUseDtokenOutput struct {
	io2.BaseResp
	EndpointTokens []EndpointToken
	Result         interface{}
}
