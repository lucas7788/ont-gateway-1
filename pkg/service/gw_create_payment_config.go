package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// CreatePaymentConfig impl
func (gw *Gateway) CreatePaymentConfig(input io.CreatePaymentConfigInput) (output io.CreatePaymentConfigOutput) {

	if !test {
		err := input.Validate()
		if err != nil {
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			return
		}
	}

	err := model.PaymentConfigManager().Insert(model.PaymentConfig(input))
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	return
}
