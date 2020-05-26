package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// ItemType ...
type ItemType int

const (
	// ItemTypeBook ...
	ItemTypeBook ItemType = iota
)

// MPEndpointGetItemMetaSchemaInput ...
type MPEndpointGetItemMetaSchemaInput struct {
	ItemType ItemType
}

// MPEndpointGetItemMetaSchemaOutput ...
type MPEndpointGetItemMetaSchemaOutput struct {
	io2.BaseResp
	ItemMetaSchema map[string]interface{}
}
