package io

import "github.com/zhiqiangxu/ont-gateway/pkg/model"

// UpdatePaymentConfigInput for input
type UpdatePaymentConfigInput model.PaymentConfig

// UpdatePaymentConfigOutput for output
type UpdatePaymentConfigOutput struct {
	BaseResp
	Exists bool `json:"exists"`
}
