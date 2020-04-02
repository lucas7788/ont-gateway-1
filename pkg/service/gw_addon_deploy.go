package service

import (
	"fmt"
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/cicd"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

const (
	defaultConfigPath = "/data/config"
	defaultStatePath  = "/data"
)

// AddonDeploy impl
func (gw *Gateway) AddonDeploy(input io.AddonDeployInput) (output io.AddonDeployOutput) {
	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	if input.ConfigPath == "" {
		input.ConfigPath = defaultConfigPath
	}
	if input.StatePath == "" {
		input.StatePath = defaultStatePath
	}

	deploymentID := fmt.Sprintf("%s:%s", input.AddonID, input.TenantID)
	ac, err := model.AddonConfigManager().Get(input.AddonID, input.TenantID)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	var config string
	if ac != nil {
		config = ac.Config
	}
	sip, err := cicd.Deploy(deploymentID, input.DockerImg, input.Spec, input.ConfigPath, config, input.StatePath)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	ad := model.AddonDeployment{AddonID: input.AddonID, TenantID: input.TenantID, SIP: sip}
	err = model.AddonDeploymentManager().Upsert(ad)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	if !input.Official {
		output.SIP = sip
	}
	return
}
