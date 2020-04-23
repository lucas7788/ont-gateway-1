package io

import (
	"fmt"

	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// GetPaymentInfoInput for input
type GetPaymentInfoInput struct {
	App        int    `json:"app"`
	PaymentID  string `json:"payment_id"`
	WithOrders bool   `json:"with_orders"`
}

// GetPaymentInfoOutput for output
type GetPaymentInfoOutput struct {
	BaseResp
	Payment       *model.Payment       `json:"payment"`
	PaymentConfig *model.PaymentConfig `json:"payment_config"`
	PaymentOrders []model.PaymentOrder `json:"payment_orders"`
}

// Validate GetPaymentInfoInput
func (input *GetPaymentInfoInput) Validate() error {
	switch {
	case input.App == 0:

		return fmt.Errorf("App empty")

	case input.PaymentID == "":

		return fmt.Errorf("PaymentID empty")

	}

	return nil
}
