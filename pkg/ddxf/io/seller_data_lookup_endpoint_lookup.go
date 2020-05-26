package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// SellerDataLookupEndpointLookupInput ...
type SellerDataLookupEndpointLookupInput struct {
	DataMetaHash string
}

// SellerDataLookupEndpointLookupOutput ...
type SellerDataLookupEndpointLookupOutput struct {
	io2.BaseResp
	DataMeta map[string]interface{}
}
