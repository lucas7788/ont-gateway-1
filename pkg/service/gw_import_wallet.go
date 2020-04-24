package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// ImportWallet impl
func (gw *Gateway) ImportWallet(input io.ImportWalletInput) (output io.ImportWalletOutput) {
	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	w := model.Wallet{Name: input.WalletName}
	err = w.SetPlainContent(input.Content)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	err = model.WalletManager().Insert(w)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	return
}
