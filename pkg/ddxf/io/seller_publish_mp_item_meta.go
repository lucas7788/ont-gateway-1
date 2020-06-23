package io

import (
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
)

// SellerPublishMPItemMetaInput ...
type SellerPublishMPItemMetaInput struct {
	ItemMeta       map[string]interface{} `bson:"itemMeta",json:"itemMeta"`
	TokenMetaHash  string                 `bson:"tokenMetaHash",json:"tokenMetaHash"`
	DataMetaHash   string                 `bson:"dataMetaHash",json:"dataMetaHash"`
	MPContractHash string                 `bson:"mpContractHash",json:"mpContractHash"`
	MPEndpoint     string                 `bson:"mpEndpoint",json:"mpEndpoint"`
	Fee            market_place_contract.Fee      `bson:"fee" json:"fee"`
	Stock          uint32                 `bson:"stock" json:"stock"`
	ExpiredDate    uint64                 `bson:"expiredDate" json:"expiredDate"`
}

type PublishParam struct {
	QrCodeId string                       `bson:"qrCodeId",json:"qrCodeId"`
	Input    SellerPublishMPItemMetaInput `bson:"input",json:"input"`
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
