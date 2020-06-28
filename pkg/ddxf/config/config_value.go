// +build !prod

package config

const (
	OperatorPrivateKey = ""

	MpUrl     = "http://127.0.0.1:" + MpPort
	SellerUrl = "http://127.0.0.1:" + SellerPort
	BuyerUrl  = "http://127.0.0.1:" + BuyerPort
)
