package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/middleware/cors"
)

const (
	buyDtoken = "/ddxf/buyer/buyDtoken"
	useDToken  = "/ddxf/buyer/useToken"

	buyDtokenQrCode     = "/onto/buyer/buyDtokenQrCode"
	qrCodeCallBack      = "/onto/buyer/qrCodeCallBack"
	getQrCodeByQrCodeId = "/onto/buyer/getQrCodeByQrCodeId/:qrCodeId"
)

var BuyerMgrAccount *ontology_go_sdk.Account

func StartBuyerServer() {
	r := gin.Default()
	r.Use(cors.Cors())
	r.POST(buyDtokenQrCode, BuyDtokenQrCodeHanler)
	r.GET(getQrCodeByQrCodeId, GetQrCodeByQrCodeIdHandler)
	r.POST(qrCodeCallBack, QrCodeCallBackHandler)
	r.POST(buyDtoken, BuyDtokenHandler)
	r.POST(useDToken, UseTokenHandler)
	err := Init()
	if err != nil {
		fmt.Println("init error:", err)
		return
	}
	private := make([]byte, 32)
	BuyerMgrAccount, err = ontology_go_sdk.NewAccountFromPrivateKey(private, signature.SHA256withECDSA)
	if err != nil {
		fmt.Println("NewAccountFromPrivateKey error:", err)
		return
	}
	go r.Run(":" + "20332")
}
