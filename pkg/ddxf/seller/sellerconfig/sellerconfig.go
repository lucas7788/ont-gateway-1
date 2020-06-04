package sellerconfig

import (
	"github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
)

type ItemMeta struct {
	ItemMetaHash string                 `bson:"itemMetaHash" json:"itemMetaHash"`
	ItemMetaData map[string]interface{} `bson:"itemMetaData" json:"itemMetaData"`
}

type SellerConfig struct {
	Version             string `json:"version"`
	ONTAuthScanProtocol string `json:"ontauth_scan_protocol"`
	QrCodeCallback      string `json:"qrcode_callback"`
	WalletName          string `json:"wallet_name"`
	ServerAccount       *ontology_go_sdk.Account
	Wallet              *ontology_go_sdk.Wallet
	Pwd                 []byte
}

var DefSellerConfig = &SellerConfig{
	Version:             "1.0",
	ONTAuthScanProtocol: "http://172.29.36.101" + config.SellerPort + "/ddxf/seller/getQrCodeDataByQrCodeId",
	QrCodeCallback:      "http://172.29.36.101" + config.SellerPort + "/ddxf/seller/qrCodeCallbackSendTx",
}
