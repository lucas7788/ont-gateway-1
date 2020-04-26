package io

import (
	"fmt"

	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// EnqueTxInput for input
type EnqueTxInput struct {
	App           int    `json:"app"`
	TxHash        string `json:"tx_hash"`
	ExpireSeconds int    `json:"expire_seconds"`
	PollAmount    bool   `json:"poll_amount"`
	Admin         bool   `json:"admin"` // not available to restful
}

// Validate impl
func (input *EnqueTxInput) Validate() (err error) {
	if !input.Admin {
		app := model.AppManager().GetByID(input.App)
		if app == nil {
			return fmt.Errorf("app not exists")
		}
		if app.TxNotifyURL == "" {
			return fmt.Errorf("tx_notify_url empty")
		}
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
