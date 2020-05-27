package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// MPEndpointQueryItemMetasInput ...
type MPEndpointQueryItemMetasInput struct {
	PageNum int64
	PageSize int64
}

// MPEndpointQueryItemMetasOutput ...
type MPEndpointQueryItemMetasOutput struct {
	io2.BaseResp
	ItemMetas []map[string]interface{}
}
