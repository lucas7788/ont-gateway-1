package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// SellerTokenLookupEndpointLookupInput ...
type SellerTokenLookupEndpointLookupInput struct {
	TokenMetaHash string
	Block         uint32
}

// SellerTokenLookupEndpointLookupOutput ...
type SellerTokenLookupEndpointLookupOutput struct {
	io2.BaseResp
	TokenMeta map[string]interface{}
}
