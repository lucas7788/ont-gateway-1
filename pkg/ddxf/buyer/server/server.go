package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/middleware/cors"
	"net/http"
)

const (
	BuyDtoken = "/ddxf/buyer/buyDtoken"
	UseDToken = "/ddxf/buyer/useToken"

	loginBuyer          = "/onto/buyer/login"
	getloginQrCode      = "/onto/buyer/getLoginQrcode/:qrCodeId"
	loginCallBack       = "/onto/buyer/loginCallBack"
	getLoginResult      = "/onto/buyer/getLoginResult"
	buyDtokenQrCode     = "/onto/buyer/buyDtokenQrCode"
	qrCodeCallBack      = "/onto/buyer/qrCodeCallBack"
	getQrCodeByQrCodeId = "/onto/buyer/getQrCodeByQrCodeId/:qrCodeId"
)

var BuyerMgrAccount *ontology_go_sdk.Account

func StartBuyerServer() error {
	r := gin.Default()
	r.Use(cors.Cors())
	r.GET(config.Ping, func(context *gin.Context) {
		context.JSON(http.StatusOK, "SUCCESS")
	})
	r.POST(loginBuyer, LoginHandler)
	r.GET(getloginQrCode, GetLoginQrCodeHandler)
	r.POST(loginCallBack, LoginCallBackHandler)
	r.GET(getLoginResult, GetLoginResultHandler)
	r.POST(buyDtokenQrCode, BuyDtokenQrCodeHanler)
	r.GET(getQrCodeByQrCodeId, GetQrCodeByQrCodeIdHandler)
	r.POST(qrCodeCallBack, QrCodeCallBackHandler)
	r.POST(BuyDtoken, BuyDtokenHandler)
	r.POST(UseDToken, UseTokenHandler)
	err := initDb()
	if err != nil {
		return err
	}
	private := make([]byte, 32)
	BuyerMgrAccount, err = ontology_go_sdk.NewAccountFromPrivateKey(private, signature.SHA256withECDSA)
	if err != nil {
		fmt.Println("NewAccountFromPrivateKey error:", err)
		return err
	}
	go r.Run(":" + config.BuyerPort)
	return nil
}
