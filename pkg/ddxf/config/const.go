package config

const (
	BuyerPort    = "20331"
	SellerPort   = "20332"
	MpPort       = "20333"
	RegistryPort = "20334"
)

const (
	PublishItemMetaUrl = "http://127.0.0.1:" + MpPort + "/ddxf/mp/publishItemMeta"
	SellerUseTokenUrl  = "http://127.0.0.1:" + SellerPort + "/ddxf/seller/useToken"
)

const (
	Key_OntId = "OntId"
	JWTAdmin  = "JWTAdmin"
	JWTAud    = "JWTAud"
)

const (
	UUID_PRE_DATAID    = "dataid_"
	UUID_PRE_QRCODE_ID = "qrcodeid_"
)
