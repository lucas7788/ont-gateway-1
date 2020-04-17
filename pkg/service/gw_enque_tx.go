package service

import (
	"net/http"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

const (
	defaultExpireSeconds = 20 * 60
)

// EnqueTx impl
func (gw *Gateway) EnqueTx(input io.EnqueTxInput) (output io.EnqueTxOutput) {

	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	expireSeconds := input.ExpireSeconds
	if expireSeconds <= 0 {
		expireSeconds = defaultExpireSeconds
	}

	now := time.Now()
	tx := model.Tx{Hash: input.TxHash, App: input.App, UpdatedAt: now, ExpireAt: now.Add(time.Second * time.Duration(expireSeconds))}
	err = model.TxManager().Upsert(tx)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	return
}
