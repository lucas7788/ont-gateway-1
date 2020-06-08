package config

const (
	BuyerPort    = "20331"
	SellerPort   = "20332"
	MpPort       = "20333"
	RegistryPort = "20334"
)

const (
	PublishItemMetaUrl = "http://127.0.0.1:" + MpPort
	SellerUrl          = "http://127.0.0.1:" + SellerPort
	BuyerUrl           = "http://127.0.0.1:" + BuyerPort
)

const (
	Key_OntId = "OntId"
	JWTAdmin  = "JWTAdmin"
	JWTAud    = "JWTAud"
)

const (
	UUID_PRE_DATAID    = "data_id_"
	UUID_PRE_QRCODE_ID = "qrcode_id_"
	UUID_RESOURCE_ID   = "resource_id_"
	UUID_PUBLISH_ID   = "publish_id_"
)
