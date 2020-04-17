package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// UpsertAddonConfig impl
func (gw *Gateway) UpsertAddonConfig(input io.UpsertAddonConfigInput) (output io.UpsertAddonConfigOutput) {
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
