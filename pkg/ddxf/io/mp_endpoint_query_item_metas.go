package io

import (
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MPEndpointQueryItemMetasInput ...
type MPEndpointQueryItemMetasInput struct {
	PageNum  int64
	PageSize int64
}

// MPEndpointQueryItemMetasOutput ...
type MPEndpointQueryItemMetasOutput struct {
	io2.BaseResp
	ItemMetas []ItemMeta
}

type ItemMeta struct {
	Id   primitive.ObjectID     `bson:"_id" json:"id"`
	Item map[string]interface{} `bson:"item" json:"item"`
}
