package io

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// SellerSaveDataMetaInput ...
type SellerSaveDataMetaInput struct {
	DataMeta     map[string]interface{} `json:"dataMeta"`
	DataMetaHash string                 `json:"dataMetaHash"`
	ResourceType byte                   `json:"resourceType"`
	Fee          param.Fee              `json:"fee"`
	Stock        uint32                 `json:"stock"`
	ExpiredDate  uint64                 `json:"expiredDate"`
	DataEndpoint string                 `json:"dataEndpoint"`
}

// SellerSaveDataMetaOutput ...
type SellerSaveDataMetaOutput struct {
	io2.BaseResp
}

type SellerSaveDataMeta struct {
	DataMeta     map[string]interface{} `bson:"dataMeta" json:"dataMeta"`
	DataMetaHash string                 `bson:"dataMetaHash" json:"dataMetaHash"`
	ResourceType byte                   `bson:"resourceType" json:"resourceType"`
	DataIds      string                 `bson:"dataIds" json:"dataIds"`
	// all below this shoudle save in dataId contract. did.
	OntId        string    `bson:"ontId" json:"ontId"`
	Fee          param.Fee `bson:"fee" json:"fee"`
	Stock        uint32    `bson:"stock" json:"stock"`
	ExpiredDate  uint64    `bson:"expiredDate" json:"expiredDate"`
	DataEndpoint string    `bson:"dataEndpoint" json:"dataEndpoint"`
}
