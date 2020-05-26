package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// MPAddRegistryInput ...
type MPAddRegistryInput struct {
	MP       string
	Endpoint string
	PubKey   string
}

// MPAddRegistryOutput ...
type MPAddRegistryOutput struct {
	io2.BaseResp
}
