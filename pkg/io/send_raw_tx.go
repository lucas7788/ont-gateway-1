package io

import (
	"fmt"

	"github.com/ontio/ontology/core/types"
)

// SendRawTxInput for input
type SendRawTxInput struct {
	Tx    *types.Transaction
	Addrs []string
}

// Validate SendRawTxInput
func (input *SendRawTxInput) Validate() (err error) {
	switch {
	case input.Tx == nil:
		err = fmt.Errorf("Tx empty")
	case len(input.Addrs) == 0:
		err = fmt.Errorf("Addrs empty")
	}
	return
}

// SendRawTxOutput for output
type SendRawTxOutput struct {
	BaseResp
	TxHash string
}
