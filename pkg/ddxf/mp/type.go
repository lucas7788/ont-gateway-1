package mp

import "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"

// Marketplace ...
type Marketplace interface {
	AddRegistry(io.MPAddRegistryInput) io.MPAddRegistryOutput
	RemoveRegistry(io.MPRemoveRegistryInput) io.MPRemoveRegistryOutput

	Endpoint() Endpoint
}

// Endpoint ...
type Endpoint interface {
	GetAuditRule(io.MPEndpointGetAuditRuleInput) io.MPEndpointGetAuditRuleOutput
	GetFee(io.MPEndpointGetFeeInput) io.MPEndpointGetFeeOutput
	GetChallengePeriod(io.MPEndpointGetChallengePeriodInput) io.MPEndpointGetChallengePeriodOutput
	GetItemMetaSchema(io.MPEndpointGetItemMetaSchemaInput) io.MPEndpointGetItemMetaSchemaOutput

	GetItemMeta(io.MPEndpointGetItemMetaInput) io.MPEndpointGetItemMetaOutput
	QueryItemMetas(io.MPEndpointQueryItemMetasInput) io.MPEndpointQueryItemMetasOutput
	PublishItemMeta(io.MPEndpointPublishItemMetaInput) io.MPEndpointPublishItemMetaOutput
}
