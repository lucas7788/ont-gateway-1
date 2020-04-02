package model

import "fmt"

func keyForAddonConfig(addonID, tenantID string) string {
	return fmt.Sprintf("ac_%s:%s", addonID, tenantID)
}

func keyForAddonDeployment(addonID, tenantID string) string {
	return fmt.Sprintf("ad_%s:%s", addonID, tenantID)
}
