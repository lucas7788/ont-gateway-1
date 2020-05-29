package io

import (
	"github.com/ontio/ontology/common"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// BuyerUseTokenInput ...
type BuyerUseTokenInput struct {
	Tx              string
	Buyer           common.Address
	Sign            string
	TokenOpEndpoint string
}

// BuyerUseTokenOutput ...
type BuyerUseTokenOutput struct {
	io2.BaseResp
	Result interface{}
}
