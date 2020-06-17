package io

import (
	"github.com/ont-bizsuite/ddxf-sdk/ddxf_contract"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// SellerSaveDataMetaInput ...
type SellerSaveDataMetaInput struct {
	DataMeta     map[string]interface{} `json:"dataMeta"`
	DataMetaHash string                 `json:"dataMetaHash"`
	ResourceType byte                   `json:"resourceType"`
	Fee          ddxf_contract.Fee      `json:"fee"`
	Stock        uint32                 `json:"stock"`
	ExpiredDate  uint64                 `json:"expiredDate"`
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
	OntId        string            `bson:"ontId" json:"ontId"`
	Fee          ddxf_contract.Fee `bson:"fee" json:"fee"`
	Stock        uint32            `bson:"stock" json:"stock"`
	ExpiredDate  uint64            `bson:"expiredDate" json:"expiredDate"`
	DataEndpoint string            `bson:"dataEndpoint" json:"dataEndpoint"`
}
