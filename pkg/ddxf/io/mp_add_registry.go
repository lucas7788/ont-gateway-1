package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// MPAddRegistryInput ...
type MPAddRegistryInput struct {
	MP       string `bson:"mp" json:"mp"`
	Endpoint string `bson:"endpoint" json:"endpoint"`
	PubKey   string `bson:"pubkey" json:"pubkey"`
}

// MPAddRegistryOutput ...
type MPAddRegistryOutput struct {
	io2.BaseResp
}
