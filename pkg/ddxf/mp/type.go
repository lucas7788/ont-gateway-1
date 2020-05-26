package mp

import "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"

// Marketplace ...
type Marketplace interface {
	PublishRegistry(io.MPPublishRegistryInput) io.MPPublishRegistryOutput
	UnPublishRegistry(io.MPUnPublishRegistryInput) io.MPUnPublishRegistryOutput

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
