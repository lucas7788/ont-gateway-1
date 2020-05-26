package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// BuyerSaveTokenAndEndpointInput ...
type BuyerSaveTokenAndEndpointInput struct {
	Tokens   map[string]interface{}
	Endpoint string
}

// BuyerSaveTokenAndEndpointOutput ...
type BuyerSaveTokenAndEndpointOutput struct {
	io2.BaseResp
}
