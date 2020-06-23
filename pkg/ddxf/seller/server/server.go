package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/middleware/jwt"
)

const (
	SaveDataMetaUrl            = "/ddxf/seller/saveDataMeta"
	GetDataIdByDataMetaHashUrl = "/ddxf/seller/getDataIdByDataMetaHash"
	SaveDataMetaArrayUrl       = "/ddxf/seller/saveDataMetaArray"
	SaveTokenMetaUrl           = "/ddxf/seller/saveTokenMeta"
	PublishMPItemMetaUrl       = "/ddxf/seller/publishMPItemMeta"
	FreezeUrl                  = "/ddxf/seller/freeze"
	UseTokenUrl                = "/ddxf/seller/useToken"
	BuyAndUseDTokenUrl         = "/ddxf/seller/buyAndUseDToken"
	PublishItemMetaUrl         = "/onto/ddxf/seller/publishItemMeta"
	getQrCodeDataByQrCodeIdUrl = "/onto/ddxf/seller/getQrCodeDataByQrCodeId"
	qrCodeCallbackSendTxUrl    = "/onto/ddxf/seller/qrCodeCallbackSendTx"
	deleteUrl = "/onto/ddxf/seller/delete"
	updateUrl = "/onto/ddxf/seller/update"
)

func StartSellerServer() {
	InitSellerImpl()
	r := gin.Default()
	r.GET(config.Ping, func(context *gin.Context) {
		context.JSON(http.StatusOK, "SUCCESS")
	})
	r.Use(jwt.JWT())
	//r.Use(middleware.JWT)
	r.POST(deleteUrl, DeleteHandler)
	r.POST(updateUrl, UpdateHandler)
	r.POST(SaveDataMetaUrl, SaveDataMetaHandler)
	r.POST(GetDataIdByDataMetaHashUrl, GetDataIdByDataMetaHashHandler)
	r.POST(SaveDataMetaArrayUrl, SaveDataMetaArrayHandler)
	r.POST(SaveTokenMetaUrl, SaveTokenMetaHandler)
	r.POST(PublishMPItemMetaUrl, PublishMPItemMetaHandler)
	r.POST(UseTokenUrl, UseTokenHandler)
	r.POST(BuyAndUseDTokenUrl, BuyAndUseDTokenHandler)
	r.POST(FreezeUrl, FreezeHandler)
	r.POST(PublishItemMetaUrl, PublishMetaHandlerOnto)
	r.POST(getQrCodeDataByQrCodeIdUrl, GetQrCodeDataByQrCodeIdHandlerOnto)
	r.POST(qrCodeCallbackSendTxUrl, GrCodeCallbackSendTxHandlerOnto)
	go r.Run(":" + config.SellerPort)
}
