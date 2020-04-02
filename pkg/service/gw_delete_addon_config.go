package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// DeleteAddonConfig impl
func (gw *Gateway) DeleteAddonConfig(input io.DeleteAddonConfigInput) (output io.DeleteAddonConfigOutput) {
	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	err = model.AddonConfigManager().Delete(input.AddonID, input.TenantID)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}
