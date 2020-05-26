package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// SellerSaveTokenMetaInput ...
type SellerSaveTokenMetaInput struct {
	DataMeta map[string]interface{}
}

// SellerSaveTokenMetaOutput ...
type SellerSaveTokenMetaOutput struct {
	io2.BaseResp
}
