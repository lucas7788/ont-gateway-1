package io

import "fmt"

// PostAddonConfigInput for input
type PostAddonConfigInput struct {
	AddonID  string `json:"addon_id"`
	TenantID string `json:"tenant_id"`
	Config   string `json:"config"`
}

// PostAddonConfigOutput for output
type PostAddonConfigOutput struct {
	BaseResp
}

// Validate PostAddonConfigInput
func (input *PostAddonConfigInput) Validate() error {
	switch {
	case input.AddonID == "":

		return fmt.Errorf("AddonID empty")

	case input.TenantID == "":

		return fmt.Errorf("TenantID empty")

	}

	return nil
}
