package io

import "fmt"

// UpsertAddonConfigInput for input
type UpsertAddonConfigInput struct {
	AddonID  string `json:"addon_id"`
	TenantID string `json:"tenant_id"`
	Config   string `json:"config"`
}

// UpsertAddonConfigOutput for output
type UpsertAddonConfigOutput struct {
	BaseResp
}

// Validate UpsertAddonConfigInput
func (input *UpsertAddonConfigInput) Validate() error {
	switch {
	case input.AddonID == "":

		return fmt.Errorf("AddonID empty")

	case input.TenantID == "":

		return fmt.Errorf("TenantID empty")

	}

	return nil
}
