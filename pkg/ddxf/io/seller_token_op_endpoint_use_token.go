package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// SellerTokenLookupEndpointUseTokenInput ...
type SellerTokenLookupEndpointUseTokenInput struct {
	Tx string
}

// SellerTokenLookupEndpointUseTokenOutput ...
type SellerTokenLookupEndpointUseTokenOutput struct {
	io2.BaseResp
	Result interface{}
}
