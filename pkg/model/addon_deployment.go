package model

import "time"

// AddonDeployment for per tenant addon deployment
type AddonDeployment struct {
	AddonID   string    `bson:"addon_id" json:"addon_id"`
	TenantID  string    `bson:"tenant_id" json:"tenant_id"`
	SIP       string    `bson:"sip" json:"sip"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
