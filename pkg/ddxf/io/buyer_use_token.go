package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// BuyerUseTokenInput ...
type BuyerUseTokenInput struct {
	Tx              string
	OnchainItemID   string
	TokenOpEndpoint string
}

// BuyerUseTokenOutput ...
type BuyerUseTokenOutput struct {
	io2.BaseResp
	Result interface{}
}
