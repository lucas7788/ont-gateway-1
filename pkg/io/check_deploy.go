package io

import "fmt"

// CheckDeployInput for input
type CheckDeployInput struct {
	AddonID  string `json:"addon_id"`
	TenantID string `json:"tenant_id"`
}

// CheckDeployOutput for output
type CheckDeployOutput struct {
	BaseResp
	OK bool `json:"ok"`
}

// Validate CheckDeployInput
func (input *CheckDeployInput) Validate() error {
	switch {
	case input.AddonID == "":

		return fmt.Errorf("AddonID empty")

	case input.TenantID == "":

		return fmt.Errorf("TenantID empty")

	}

	return nil
}
