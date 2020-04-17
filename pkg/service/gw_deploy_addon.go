package service

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"crypto/md5"

	"github.com/zhiqiangxu/ont-gateway/pkg/cicd"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

const (
	defaultConfigPath = "/appconfig/config.json"
	defaultStatePath  = "/data"
)

// the resultant length should be <= 64 for cicd
func toDeploymentID(addonID, tenantID string) string {
	longID := fmt.Sprintf("%s:%s", addonID, tenantID)
	sum := md5.Sum([]byte(longID))
	return hex.EncodeToString(sum[:])
}

// DeployAddon impl
func (gw *Gateway) DeployAddon(input io.DeployAddonInput) (output io.DeployAddonOutput) {
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

	deploymentID := toDeploymentID(input.AddonID, input.TenantID)
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

	ad := model.AddonDeployment{AddonID: input.AddonID, TenantID: input.TenantID, SIP: sip, UpdatedAt: time.Now()}
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
