package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// MPPublishRegistryInput ...
type MPPublishRegistryInput struct {
	MP       string
	Endpoint string
	PubKey   string
}

// MPPublishRegistryOutput ...
type MPPublishRegistryOutput struct {
	io2.BaseResp
}
