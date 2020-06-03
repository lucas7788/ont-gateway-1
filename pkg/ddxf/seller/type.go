package seller

import "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"

// Seller ...
type Seller interface {
	SaveDataMeta(io.SellerSaveDataMetaInput) io.SellerSaveDataMetaOutput
	SaveTokenMeta(io.SellerSaveTokenMetaInput) io.SellerSaveTokenMetaOutput
	PublishMPItemMeta(io.SellerPublishMPItemMetaInput) io.SellerPublishMPItemMetaOutput

	
	DataLookupEndpoint() DataLookupEndpoint
	TokenLookupEndpoint() TokenLookupEndpoint
	TokenOpEndpoint() TokenOpEndpoint
}

// DataLookupEndpoint ...
type DataLookupEndpoint interface {
	Lookup(io.SellerDataLookupEndpointLookupInput) io.SellerDataLookupEndpointLookupOutput
}

// TokenLookupEndpoint ...
type TokenLookupEndpoint interface {
	Lookup(io.SellerTokenLookupEndpointLookupInput) io.SellerTokenLookupEndpointLookupOutput
}

// TokenOpEndpoint ...
type TokenOpEndpoint interface {
	UseToken(io.SellerTokenLookupEndpointUseTokenInput) io.SellerTokenLookupEndpointUseTokenOutput
}
