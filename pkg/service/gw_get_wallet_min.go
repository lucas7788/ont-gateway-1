package service

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// GetWalletMin impl
func (gw *Gateway) GetWalletMin(input io.GetWalletMinInput) (output io.GetWalletMinOutput) {
	out := gw.GetWallet(io.GetWalletInput(input))
	if out.Error() != nil {
		output.SetBase(out.Code, out.Msg)
		return
	}

	output.Content = out.Content
	output.Exists = out.Exists

	return
}
