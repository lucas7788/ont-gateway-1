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

func AddEndpointHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	param := io.MPAddRegistryInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	output := AddEndpointService(io.RegistryAddEndpointInput(param))
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, output)
}

func RemoveEndpointHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[RemoveEndpointHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	param := io.MPRemoveRegistryInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[RemoveEndpointHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	output := RemoveEndpointService(io.RegistryRemoveEndpointInput(param))
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusInternalServerError, output)
		return
	}
	ctx.JSON(output.Code, output)
}

func QueryEndpointHandler(ctx *gin.Context) {
	output := QueryEndpointsService(io.RegistryQueryEndpointsInput{})
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusInternalServerError, output)
		return
	}
	ctx.JSON(http.StatusOK, output)
}
