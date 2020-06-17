package server

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

type ItemMeta struct {
	ItemMetaHash string                 `bson:"itemMetaHash" json:"itemMetaHash"`
	ItemMetaData map[string]interface{} `bson:"itemMetaData" json:"itemMetaData"`
}

type PublishForOpenKgParam struct {
	SellerSaveDataMetaInput        io.SellerSaveDataMetaInput
	SellerSaveTokenMetaInput       io.SellerSaveTokenMetaInput
	MPEndpointPublishItemMetaInput io.MPEndpointPublishItemMetaInput
}

type OpenKgRes struct {
	io2.BaseResp
}

type FreezeParam struct {
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
