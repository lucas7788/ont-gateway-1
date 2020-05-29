package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func BuyDtokenHandler(ctx *gin.Context) {
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
	output := BuyDTokenService(param)
	ctx.JSON(output.Code, output)
}

func UseTokenHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[UseTokenHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	param := io.BuyerUseTokenInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[UseTokenHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	output := UseTokenService(param)
	ctx.JSON(output.Code, output)
}
