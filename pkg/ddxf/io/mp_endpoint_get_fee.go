package io

import (
	"github.com/ont-bizsuite/ddxf-sdk/ddxf_contract"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// MPEndpointGetFeeInput ...
type MPEndpointGetFeeInput struct {
}

// MPEndpointGetFeeOutput ...
type MPEndpointGetFeeOutput struct {
	io.BaseResp
	ddxf_contract.Fee
}
