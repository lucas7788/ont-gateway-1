package io

import (
	"fmt"
)

// SendTxInput for input
type SendTxInput struct {
	SignedTx string
	Addrs    []string
}

// Validate SendTxInput
func (input *SendTxInput) Validate() (err error) {
	switch {
	case input.SignedTx == "":
		err = fmt.Errorf("SignedTx empty")
	case len(input.Addrs) == 0:
		err = fmt.Errorf("Addrs empty")
	}
	return
}

// SendTxOutput for output
type SendTxOutput struct {
	BaseResp
	TxHash string
}
