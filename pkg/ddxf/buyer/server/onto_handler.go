package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func LoginHandler(ctx *gin.Context) {
	res := LoginService()
	ctx.JSON(http.StatusOK, ResponseSuccess(res))
}

func GetLoginQrCodeHandler(ctx *gin.Context) {
	qrCodeId := ctx.Param("qrCodeId")
	code, err := GetLoginQrCodeService(qrCodeId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		ctx.JSON(http.StatusOK, ResponseSuccess(code))
	}
}

func LoginCallBackHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[LoginCallBackHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	param := QrCodeCallBackParam{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[LoginCallBackHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	err = LoginCallBackService(param)
	if err != nil {
		instance.Logger().Error("[LoginCallBackHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	ctx.JSON(http.StatusOK, ResponseSuccess("SUCCESS"))
}

func GetLoginResultHandler(ctx *gin.Context) {
	qrCodeId := ctx.Param("qrCodeId")
	status, err := GetLoginResultService(qrCodeId)
	if err != nil {
		instance.Logger().Error("[GetLoginResultHandler] param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	ctx.JSON(http.StatusOK, ResponseSuccess(status))
}

func BuyDtokenQrCodeHanler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[BuyDtokenHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	param := BuyerBuyDtokenQrCodeInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[BuyDtokenHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	output, err := BuyDtokenQrCodeService(param)
	if err != nil {
		instance.Logger().Error("[BuyDtokenHandler] BuyDtokenService error:", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	ctx.JSON(http.StatusOK, ResponseSuccess(output))
}

func GetQrCodeByQrCodeIdHandler(ctx *gin.Context) {
	qrCodeId := ctx.Param("qrCodeId")
	code, err := GetQrCodeByQrCodeIdService(qrCodeId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		ctx.JSON(http.StatusOK, ResponseSuccess(code))
	}
}

func QrCodeCallBackHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[QrCodeCallBackHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	param := QrCodeCallBackParam{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[QrCodeCallBackHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	output, err := QrCodeCallBackService(param)
	if err != nil {
		instance.Logger().Error("[QrCodeCallBackHandler] QrCodeCallBackService error:", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		ctx.JSON(http.StatusOK, ResponseSuccess(output))
	}
}
