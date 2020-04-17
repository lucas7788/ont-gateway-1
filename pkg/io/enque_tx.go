package io

import "fmt"

// EnqueTxInput for input
type EnqueTxInput struct {
	App           int    `json:"app"`
	TxHash        string `json:"tx_hash"`
	ExpireSeconds int    `json:"expire_seconds"`
}

// Validate impl
func (input *EnqueTxInput) Validate() (err error) {
	if input.App == 0 {
		err = fmt.Errorf("empty app")
		return
	}
	if input.TxHash == "" {
		err = fmt.Errorf("empty TxHash")
		return
	}
	return
}

// EnqueTxOutput for output
type EnqueTxOutput struct {
	BaseResp
}
