package io

import "fmt"

// DeployAddonInput for input
type DeployAddonInput struct {
	AddonID    string `json:"addon_id"`
	TenantID   string `json:"tenant_id"`
	Official   bool   `json:"official"`
	DockerImg  string `json:"docker_img"`
	Spec       string `json:"spec"`
	ConfigPath string `json:"config_path"`
	StatePath  string `json:"state_path"`
}

// DeployAddonOutput for output
type DeployAddonOutput struct {
	BaseResp
	SIP string `json:"sip"`
}

// Validate DeployAddonInput
func (input *DeployAddonInput) Validate() error {
	switch {
	case input.AddonID == "":

		return fmt.Errorf("AddonID empty")

	case input.TenantID == "":

		return fmt.Errorf("TenantID empty")

	case input.DockerImg == "":

		return fmt.Errorf("DockerImg empty")

	}

	return nil
}
