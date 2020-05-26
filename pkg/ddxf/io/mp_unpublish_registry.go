package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// MPUnPublishRegistryInput ...
type MPUnPublishRegistryInput struct {
	MP   string
	Sign string
}

// MPUnPublishRegistryOutput ...
type MPUnPublishRegistryOutput struct {
	io2.BaseResp
}
