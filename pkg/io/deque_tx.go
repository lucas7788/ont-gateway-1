package io

import "fmt"

// DequeTxInput for input
type DequeTxInput struct {
	App    int    `json:"app"`
	TxHash string `json:"tx_hash"`
	Admin  bool   `json:"admin"` // not available to restful
}

// Validate impl
func (input *DequeTxInput) Validate() (err error) {
	if !input.Admin {
		if input.App == 0 {
			err = fmt.Errorf("empty app")
			return
		}
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
