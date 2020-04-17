package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// CreateApp impl
func (gw *Gateway) CreateApp(input io.CreateAppInput) (output io.CreateAppOutput) {

	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	err = model.AppManager().Insert(model.App(input))
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	return
}
