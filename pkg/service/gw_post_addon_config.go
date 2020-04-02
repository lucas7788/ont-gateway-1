package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// PostAddonConfig impl
func (gw *Gateway) PostAddonConfig(input io.PostAddonConfigInput) (output io.PostAddonConfigOutput) {
	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	ac := model.AddonConfig{AddonID: input.AddonID, TenantID: input.TenantID, Config: input.Config}
	err = model.AddonConfigManager().Upsert(ac)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}
