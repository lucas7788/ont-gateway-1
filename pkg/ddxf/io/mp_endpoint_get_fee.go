package io

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
)

// MPEndpointGetFeeInput ...
type MPEndpointGetFeeInput struct {
}

// MPEndpointGetFeeOutput ...
type MPEndpointGetFeeOutput struct {
	io.BaseResp
	market_place_contract.Fee
}
