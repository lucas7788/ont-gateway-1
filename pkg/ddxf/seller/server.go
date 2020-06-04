package seller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/restful"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/service"
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
	service.InitSellerImpl()
	r := gin.Default()
	//r.Use(middleware.JWT)
	r.POST(SaveDataMetaUrl, restful.SaveDataMetaHandle)
	r.POST(SaveTokenMetaUrl, restful.SaveTokenMetaHandle)
	r.POST(PublishMPItemMetaUrl, restful.PublishMPItemMetaHandle)
	r.POST(UseTokenUrl, restful.UseTokenHandler)
	r.POST(PublishItemMetaUrl, restful.PublishMetaHandler)
	r.POST(getQrCodeDataByQrCodeIdUrl, restful.GetQrCodeDataByQrCodeId)
	r.POST(qrCodeCallbackSendTxUrl, restful.GrCodeCallbackSendTx)
	go r.Run(":" + config.SellerPort)
}
