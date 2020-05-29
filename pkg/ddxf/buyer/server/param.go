package server

import "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"

type QrCodeAndEndpoint struct {
	Code          qrCode.QrCode `bson:"qrCode",json:"qrCode"`
	TokenEndpoint string        `bson:"tokenEndpoint",json:"tokenEndpoint"`
}

type BuyerBuyDtokenQrCodeInput struct {
	OnchainItemId   string
	N               int
	Buyer           string
	TokenOpEndpoint string
}

type QrCodeCallBackParam struct {
	Signer    string    `json:"signer"`
	SignedTx  string    `json:"signedTx"`
	ExtraData ExtraData `json:"extraData"`
}
type ExtraData struct {
	Id        string `json:"id"`
	PublicKey string `json:"publickey"`
	OntId     string `json:"ontId"`
}

func ResponseFailedOnto(errCode int64, err error) map[string]interface{} {
	return map[string]interface{}{
		"result":  err,
		"error":   errCode,
		"desc":    err.Error(),
		"version": "1.0",
	}
}

func ResponseSuccess(result interface{}) map[string]interface{} {
	return map[string]interface{}{
		"result":  result,
		"error":   0,
		"desc":    "SUCCESS",
		"version": "1.0",
	}
}
