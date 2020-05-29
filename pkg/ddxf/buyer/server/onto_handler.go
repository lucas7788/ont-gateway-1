package server

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

func BuyDtokenQrCodeHanler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[BuyDtokenHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	param := io.BuyerBuyDtokenInput{}
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
	ctx.JSON(0, ResponseSuccess(output))
}

func GetQrCodeByQrCodeIdHandler(ctx *gin.Context) {
	qrCodeId := ctx.Param("qrCodeId")
	code, err := GetQrCodeByQrCodeIdService(qrCodeId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		ctx.JSON(0, ResponseSuccess(code))
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
		ctx.JSON(0, ResponseSuccess(output))
	}
}

