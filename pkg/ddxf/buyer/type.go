package buyer

import "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"

// Buyer ...
type Buyer interface {
	BuyDtoken(io.BuyerBuyDtokenInput) io.BuyerBuyDtokenOutput
	SaveTokenAndEndpoint(io.BuyerSaveTokenAndEndpointInput) io.BuyerSaveTokenAndEndpointOutput

	UseToken(io.BuyerUseTokenInput) io.BuyerUseTokenOutput
}
