package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
)

const (
	SaveDataMetaUrl            = "/ddxf/seller/saveDataMeta"
	SaveTokenMetaUrl           = "/ddxf/seller/saveTokenMeta"
	PublishMPItemMetaUrl       = "/ddxf/seller/publishMPItemMeta"
	UseTokenUrl                = "/ddxf/seller/useToken"
	PublishItemMetaUrl         = "/ddxf/seller/publishItemMeta"
	getQrCodeDataByQrCodeIdUrl = "ddxf/seller/getQrCodeDataByQrCodeId"
	qrCodeCallbackSendTxUrl    = "ddxf/seller/qrCodeCallbackSendTx"
)

func StartSellerServer() {
	InitSellerImpl()
	r := gin.Default()
	//r.Use(middleware.JWT)
	r.POST(SaveDataMetaUrl, SaveDataMetaHandler)
	r.POST(SaveTokenMetaUrl, SaveTokenMetaHandler)
	r.POST(PublishMPItemMetaUrl, PublishMPItemMetaHandler)
	r.POST(UseTokenUrl, UseTokenHandler)
	r.POST(PublishItemMetaUrl, PublishMetaHandler)
	r.POST(getQrCodeDataByQrCodeIdUrl, GetQrCodeDataByQrCodeIdHandler)
	r.POST(qrCodeCallbackSendTxUrl, GrCodeCallbackSendTxHandler)
	go r.Run(":" + config.SellerPort)
}
