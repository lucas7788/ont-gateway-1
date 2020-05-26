package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// SellerSaveDataMetaInput ...
type SellerSaveDataMetaInput struct {
	DataMeta map[string]interface{}
}

// SellerSaveDataMetaOutput ...
type SellerSaveDataMetaOutput struct {
	io2.BaseResp
}
