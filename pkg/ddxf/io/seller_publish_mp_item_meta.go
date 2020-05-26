package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// SellerPublishMPItemMetaInput ...
type SellerPublishMPItemMetaInput struct {
	DataMeta   map[string]interface{}
	MPEndpoint string
}

// SellerPublishMPItemMetaOutput ...
type SellerPublishMPItemMetaOutput struct {
	io2.BaseResp
}
