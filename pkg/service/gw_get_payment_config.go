package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// GetPaymentConfig impl
func (gw *Gateway) GetPaymentConfig(input io.GetPaymentConfigInput) (output io.GetPaymentConfigOutput) {

	config, err := model.PaymentConfigManager().Get(input.App, input.PaymentConfigID)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	if config != nil {
		output.Exists = true
		output.PaymentConfig = *config
	}

	return
}
