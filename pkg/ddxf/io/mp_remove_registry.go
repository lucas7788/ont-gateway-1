package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// MPRemoveRegistryInput ...
type MPRemoveRegistryInput struct {
	MP   string
	Sign []byte
}

// MPRemoveRegistryOutput ...
type MPRemoveRegistryOutput struct {
	io2.BaseResp
}
