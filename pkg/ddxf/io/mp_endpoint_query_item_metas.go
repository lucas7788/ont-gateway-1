package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// MPEndpointQueryItemMetasInput ...
type MPEndpointQueryItemMetasInput struct {
	Text string
}

// MPEndpointQueryItemMetasOutput ...
type MPEndpointQueryItemMetasOutput struct {
	io2.BaseResp
	ItemMetas []map[string]interface{}
}
