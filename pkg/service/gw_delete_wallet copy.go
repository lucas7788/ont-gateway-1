package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// DeleteWallet impl
func (gw *Gateway) DeleteWallet(input io.DeleteWalletInput) (output io.DeleteWalletOutput) {

	err := model.WalletManager().DeleteOne(input.WalletName)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	return
}
