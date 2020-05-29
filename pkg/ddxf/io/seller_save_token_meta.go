package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// SellerSaveTokenMetaInput ...
type SellerSaveTokenMetaInput struct {
	TokenMeta     map[string]interface{} `json:"tokenMeta"`
	TokenMetaHash string                 `json:"tokenMetaHash"`
	DataMetaHash  string                 `json:"dataMetaHash"`
	TokenEndpoint string                 `json:"tokenEndpoint"`
}

// SellerSaveTokenMetaOutput ...
type SellerSaveTokenMetaOutput struct {
	io2.BaseResp
}

type SellerSaveTokenMeta struct {
	TokenMeta     map[string]interface{} `bson:"tokenMeta" json:"tokenMeta"`
	TokenMetaHash string                 `bson:"tokenMetaHash" json:"tokenMetaHash"`
	DataMetaHash  string                 `bson:"dataMetaHash" json:"dataMetaHash"`
	TokenEndpoint string                 `bson:"tokenEndpoint" json:"tokenEndpoint"`
	OntId         string                 `bson:"ontId" json:"ontId"`
}
