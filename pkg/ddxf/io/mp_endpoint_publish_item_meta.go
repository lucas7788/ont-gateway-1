package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// MPEndpointPublishItemMetaInput ...
type MPEndpointPublishItemMetaInput struct {
	ItemMeta     map[string]interface{}
	SignedDDXFTx string
}

// MPEndpointPublishItemMetaOutput ...
type MPEndpointPublishItemMetaOutput struct {
	io2.BaseResp
}
