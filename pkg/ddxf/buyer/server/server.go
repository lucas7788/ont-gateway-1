package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/middleware/cors"
)

const (
	buyDtoken = "/ddxf/buyer/buyDtoken"
	useToken  = "/ddxf/buyer/useToken"

	buyDtokenQrCode     = "/onto/buyer/buyDtokenQrCode"
	qrCodeCallBack      = "/onto/buyer/qrCodeCallBack"
	getQrCodeByQrCodeId = "/onto/buyer/getQrCodeByQrCodeId/:qrCodeId"
)

func StartBuyerServer() {
	r := gin.Default()
	r.Use(cors.Cors())
	r.POST(buyDtokenQrCode, BuyDtokenQrCodeHanler)
	r.GET(getQrCodeByQrCodeId, GetQrCodeByQrCodeIdHandler)
	r.POST(qrCodeCallBack, QrCodeCallBackHandler)
	r.POST(buyDtoken, BuyDtokenHandler)
	r.POST(useToken, UseTokenHandler)
	Init()
	go r.Run(":" + "20332")
}
