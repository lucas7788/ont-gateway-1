package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// GetAddonConfig impl
func (gw *Gateway) GetAddonConfig(input io.GetAddonConfigInput) (output io.GetAddonConfigOutput) {
	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	ac, err := model.AddonConfigManager().Get(input.AddonID, input.TenantID)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	if ac == nil {
		return
	}

	output.Exists = true
	output.Config = ac.Config
	return
}
