package io

import "fmt"

// DequeTxInput for input
type DequeTxInput struct {
	App    int    `json:"app"`
	TxHash string `json:"tx_hash"`
}

// Validate impl
func (input *DequeTxInput) Validate() (err error) {
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

// DequeTxOutput for output
type DequeTxOutput struct {
	BaseResp
	Exists bool `json:"exists"`
}
