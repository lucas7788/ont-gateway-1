package io

import (
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// SellerSaveDataMetaInput ...
type SellerSaveDataMetaInput struct {
	DataMeta     map[string]interface{} `json:"dataMeta"`
	DataMetaHash string                 `json:"dataMetaHash"`
	ResourceType byte                   `json:"resourceType"`
	DataEndpoint string                 `json:"dataEndpoint"`
	DataHash     string                 `json:"dataHash"`
	SignedTx     string                 `json:"signedTx"`
	DataId       string                 `json:"dataId"`
}

// SellerSaveDataMetaOutput ...
type SellerSaveDataMetaOutput struct {
	io2.BaseResp
	DataId string
}

type SellerSaveDataMeta struct {
	DataMeta     map[string]interface{} `bson:"dataMeta" json:"dataMeta"`
	DataMetaHash string                 `bson:"dataMetaHash" json:"dataMetaHash"`
	ResourceType byte                   `bson:"resourceType" json:"resourceType"`
	DataId       string                 `bson:"dataId" json:"dataId"`
	// all below this shoudle save in dataId contract. did.
	OntId        string `bson:"ontId" json:"ontId"`
	DataEndpoint string `bson:"dataEndpoint" json:"dataEndpoint"`
}

type SellerSaveDataMetaArrayInput struct {
	DataMetaOneArray []DataMetaOne `json:"dataMetaOne"`
	SignedTx         string        `json:"signedTx"`
}

type DataMetaOne struct {
	DataMeta     map[string]interface{} `json:"dataMeta"`
	DataMetaHash string                 `json:"dataMetaHash"`
	ResourceType byte                   `json:"resourceType"`
	DataEndpoint string                 `json:"dataEndpoint"`
	DataHash     string                 `json:"dataHash"`
	DataId       string                 `json:"dataId"`
}
type SellerSaveDataMetaArrayOutput struct {
	io2.BaseResp
	DataId string
}
