// +build prod

package config

const (
	OperatorPrivateKey = ""

	MpUrl     = "http://ddxf-component-mp-service:" + MpPort
	SellerUrl = "http://ddxf-component-seller-service:" + SellerPort
	BuyerUrl  = "http://ddxf-component-buyer-service:" + BuyerPort
)
