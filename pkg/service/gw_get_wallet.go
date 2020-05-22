package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// GetWallet impl
func (gw *Gateway) GetWallet(input io.GetWalletInput) (output io.GetWalletOutput) {
	w, err := model.WalletManager().GetOne(input.WalletName)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	if w == nil {
		return
	}

	psw, content, err := w.GetPlain()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	output.Exists = true
	output.PSW = psw
	output.Content = content

	return
}
