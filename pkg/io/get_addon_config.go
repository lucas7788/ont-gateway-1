package io

import "fmt"

// GetAddonConfigInput for input
type GetAddonConfigInput struct {
	AddonID  string `json:"addon_id"`
	TenantID string `json:"tenant_id"`
}

// GetAddonConfigOutput for output
type GetAddonConfigOutput struct {
	BaseResp
	Exists bool   `json:"exists"`
	Config string `json:"config"`
}

// Validate GetAddonConfigInput
func (input *GetAddonConfigInput) Validate() error {
	switch {
	case input.AddonID == "":

		return fmt.Errorf("AddonID empty")

	case input.TenantID == "":

		return fmt.Errorf("TenantID empty")

	}

	return nil
}
