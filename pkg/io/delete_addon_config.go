package io

import "fmt"

// DeleteAddonConfigInput for input
type DeleteAddonConfigInput struct {
	AddonID  string `json:"addon_id"`
	TenantID string `json:"tenant_id"`
}

// DeleteAddonConfigOutput for output
type DeleteAddonConfigOutput struct {
	BaseResp
}

// Validate DeleteAddonConfigInput
func (input *DeleteAddonConfigInput) Validate() error {
	switch {
	case input.AddonID == "":

		return fmt.Errorf("AddonID empty")

	case input.TenantID == "":

		return fmt.Errorf("TenantID empty")

	}

	return nil
}
