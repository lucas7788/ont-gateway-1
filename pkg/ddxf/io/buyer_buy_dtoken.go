package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// BuyerBuyDtokenInput ...
type BuyerBuyDtokenInput struct {
	Tx string
}

// BuyerBuyDtokenOutput ...
type BuyerBuyDtokenOutput struct {
	io2.BaseResp
	Tokens   map[string]interface{}
	Endpoint string
}
