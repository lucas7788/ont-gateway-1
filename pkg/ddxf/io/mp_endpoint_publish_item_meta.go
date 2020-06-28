package io

import (
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// MPEndpointPublishItemMetaInput ...
type MPEndpointPublishItemMetaInput struct {
	SignedDDXFTx string `bson:"signed_ddxf_tx" json:"signed_ddxf_tx"`
	ItemMeta     PublishItemMeta
	MPEndpoint   string
}

type PublishItemMeta struct {
	OnchainItemID string                 `bson:"onchain_item_id" json:"onchain_item_id"` //resource_id
	ItemMeta      map[string]interface{} `bson:"item_meta" json:"item_meta"`
}

// MPEndpointPublishItemMetaOutput ...
type MPEndpointPublishItemMetaOutput struct {
	io2.BaseResp
}
