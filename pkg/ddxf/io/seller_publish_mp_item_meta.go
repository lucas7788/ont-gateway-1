package io

import (
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// SellerPublishMPItemMetaInput ...
type SellerPublishMPItemMetaInput struct {
	ItemMeta       map[string]interface{} `bson:"itemMeta",json:"itemMeta"`
	TokenMetaHash  string                 `bson:"tokenMetaHash",json:"tokenMetaHash"`
	DataMetaHash   string                 `bson:"dataMetaHash",json:"dataMetaHash"`
	MPContractHash string                 `bson:"mpContractHash",json:"mpContractHash"`
	MPEndpoint     string                 `bson:"mpEndpoint",json:"mpEndpoint"`
}

//type TokenTemplate2 struct {
//	dataid    string
//	tokenHash []string
//}
//
//[]map[string]interface{}
//tokenHash'

// SellerPublishMPItemMetaOutput ...
type SellerPublishMPItemMetaOutput struct {
	io2.BaseResp
}
