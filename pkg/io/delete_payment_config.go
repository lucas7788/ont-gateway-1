package io

// DeletePaymentConfigInput for input
type DeletePaymentConfigInput struct {
	App             int    `json:"app"`
	PaymentConfigID string `json:"payment_config_id"`
}

// DeletePaymentConfigOutput for output
type DeletePaymentConfigOutput struct {
	BaseResp
	Exists bool `json:"exists"`
}
