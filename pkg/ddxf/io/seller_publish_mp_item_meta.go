package io

import (
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// SellerPublishMPItemMetaInput ...
type SellerPublishMPItemMetaInput struct {
	ItemMeta       map[string]interface{}
	TokenMetaHash  string
	DataMetaHash   string
	MPContractHash string
	MPEndpoint     string
}

// SellerPublishMPItemMetaOutput ...
type SellerPublishMPItemMetaOutput struct {
	io2.BaseResp
}
