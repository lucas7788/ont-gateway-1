package io

import (
	"fmt"

	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// CreatePaymentOrderInput for input
type 
CreatePaymentOrderInput struct {
	App             int             `bson:"app" json:"app"`
	PaymentConfigID string          `json:"payment_config_id"`
	PaymentID       string          `json:"payment_id"`
	OrderID         string          `json:"order_id"`
	PayPeriod       model.PayPeriod `json:"pay_period"`
	PayMethod       model.PayMethod `json:"pay_method"`
	Amount          int             `json:"amount"`
	CoinType        model.CoinType  `json:"coin_type"`
	OrderInfo       string          `json:"order_info"`
}

// CreatePaymentOrderOutput for output
type CreatePaymentOrderOutput struct {
	BaseResp
	Balance int `json:"balance"`
}

// Validate CreatePaymentOrderInput
func (input *CreatePaymentOrderInput) Validate() error {
	_, exists := model.AppManager().GetApp(input.App)
	if !exists {
		return fmt.Errorf("app not exists")
	}

	switch {
	case input.PaymentConfigID == "":
		return fmt.Errorf("payment_config_id empty")
	case input.PaymentID == "":
		return fmt.Errorf("payment_id empty")
	case input.OrderID == "":
		return fmt.Errorf("order_id empty")
	case input.PayPeriod == 0:
		return fmt.Errorf("pay_period empty")
	default:
		return model.VerifyOrderInfo(input.Amount, input.CoinType, input.OrderInfo)
	}
}

// ToPaymentOrder converts CreatePaymentOrderInput to model.PaymentOrder
func (input *CreatePaymentOrderInput) ToPaymentOrder() model.PaymentOrder {
	return model.PaymentOrder{
		App:       input.App,
		PaymentID: input.PaymentID,
		Amount:    input.Amount,
		CoinType:  input.CoinType,
		OrderInfo: input.OrderInfo,
	}
}
