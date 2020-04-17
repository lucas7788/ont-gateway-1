package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// DequeTx impl
func (gw *Gateway) DequeTx(input io.DequeTxInput) (output io.DequeTxOutput) {

	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	exists, err := model.TxManager().Delete(input.TxHash)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	output.Exists = exists

	return
}
