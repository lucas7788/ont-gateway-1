package model

// AddonConfig for per tenant addon config
type AddonConfig struct {
	AddonID  string `bson:"addon_id" json:"addon_id"`
	TenantID string `bson:"tenant_id" json:"tenant_id"`
	Config   string `bson:"config" json:"config"`
}
