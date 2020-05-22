package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// ListApp impl
func (gw *Gateway) ListApp() (output io.ListAppOutput) {

	apps, err := model.AppManager().GetAllFromDB()
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	output.Apps = apps

	return
}
