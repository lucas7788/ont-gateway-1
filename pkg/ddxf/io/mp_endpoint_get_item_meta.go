package io

// MPEndpointGetItemMetaInput ...
type MPEndpointGetItemMetaInput struct {
	ItemMetaID string
}

// MPEndpointGetItemMetaOutput ...
type MPEndpointGetItemMetaOutput struct {
	ItemMeta      map[string]interface{}
	OnchainItemID string
}
