package service

import (
	"fmt"
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/cicd"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// CheckDeploy impl
func (gw *Gateway) CheckDeploy(input io.CheckDeployInput) (output io.CheckDeployOutput) {
	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	ad, exists, err := model.AddonDeploymentManager().Get(input.AddonID, input.TenantID)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	if !exists {
		return
	}

	url := fmt.Sprintf("http://%s/health", ad.SIP)
	ok, err := cicd.Check(url)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	output.OK = ok

	return
}
