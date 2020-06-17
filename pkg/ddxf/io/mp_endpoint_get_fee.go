package io

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/ont-bizsuite/ddxf-sdk/ddxf_contract"
)

// MPEndpointGetFeeInput ...
type MPEndpointGetFeeInput struct {
}

// MPEndpointGetFeeOutput ...
type MPEndpointGetFeeOutput struct {
	io.BaseResp
	ddxf_contract.Fee
}
