package io

import (
	"fmt"

	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// CreatePaymentConfigInput for input
type CreatePaymentConfigInput model.PaymentConfig

// CreatePaymentConfigOutput for output
type CreatePaymentConfigOutput struct {
	BaseResp
}

// Validate CreatePaymentConfigInput
func (input *CreatePaymentConfigInput) Validate() error {
	app, exists := model.AppManager().GetApp(input.App)
	if !exists {
		return fmt.Errorf("app not exists")
	}
	if app.PaymentNotifyURL == "" {
		return fmt.Errorf("payment_notify_url empty")
	}

	switch {
	case input.PaymentConfigID == "":
		return fmt.Errorf("payment_config_id empty")
	case len(input.PeriodOptions) != len(input.AmountOptions):
		return fmt.Errorf("period options not match amount options")
	default:
		return nil
	}
}
