package io

import "github.com/zhiqiangxu/ont-gateway/pkg/model"

// GetPaymentConfigInput for input
type GetPaymentConfigInput struct {
	App             int    `json:"app"`
	PaymentConfigID string `json:"payment_config_id"`
}

// GetPaymentConfigOutput for output
type GetPaymentConfigOutput struct {
	BaseResp
	Exists        bool                `json:"exists"`
	PaymentConfig model.PaymentConfig `json:"payment_config"`
}
