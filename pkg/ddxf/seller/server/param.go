package server

import (
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

type ItemMeta struct {
	ItemMetaHash string                 `bson:"itemMetaHash" json:"itemMetaHash"`
	ItemMetaData map[string]interface{} `bson:"itemMetaData" json:"itemMetaData"`
}

type FreezeOutput struct {
	io2.BaseResp
}

type DeleteParam struct {
	SignedTx string
}

type GetDataIdParam struct {
	DataMetaHashArray []string `json:"dataMetaHashArray"`
}

type GetDataIdRes struct {
	DataIdAndDataMetaHashArray []*DataIdAndDataMetaHash `json:"dataMetaHashArray"`
}

type DataIdAndDataMetaHash struct {
	DataId       string
	DataMetaHash string
}

type DeleteInput struct {
	SignedTx string
}

type DeleteOutput struct {
	io2.BaseResp
	Result interface{}
}

type UpdateInput struct {
	SignedTx string
}

type UpdateOutput struct {
	io2.BaseResp
}
