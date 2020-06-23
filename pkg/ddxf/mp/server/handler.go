package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func AddRegistryHandler(ctx *gin.Context) {
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
	output := AddRegistryService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, output)
}

func RemoveRegistryHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	param := io.MPRemoveRegistryInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	output := RemoveRegistryService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusInternalServerError, output)
		return
	}
	ctx.JSON(http.StatusOK, output)
}
func PublishItemMetaHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	param := io.MPEndpointPublishItemMetaInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println("param:", param)
	output := PublishItemMetaService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusInternalServerError, output)
		return
	}
	ctx.JSON(http.StatusOK, output)
}

func DeleteHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[DeleteHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	param := DeleteInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[DeleteHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	output := DeleteService(param)
	if output.Code == 0 {
		output.Code = http.StatusOK
	}
	ctx.JSON(output.Code, output)
}

func UpdateHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[UpdateHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	param := UpdateInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[UpdateHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	output := UpdateService(param)
	if output.Code == 0 {
		output.Code = http.StatusOK
	}
	ctx.JSON(output.Code, output)
}

func GetAuditRuleHandler(ctx *gin.Context) {
	output := GetAuditRuleService(io.MPEndpointGetAuditRuleInput{})
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusInternalServerError, output)
		return
	}
	ctx.JSON(http.StatusOK, output)
}
func GetFeeHandler(ctx *gin.Context) {
	output := GetFeeService(io.MPEndpointGetFeeInput{})
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusInternalServerError, output)
		return
	}
	ctx.JSON(http.StatusOK, output)
}
func GetChallengePeriodHandler(ctx *gin.Context) {
	output := GetChallengePeriodService(io.MPEndpointGetChallengePeriodInput{})
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusBadRequest, output)
		return
	}
	ctx.JSON(http.StatusOK, output)
}
func GetItemMetaSchemaHandler(ctx *gin.Context) {
	output := GetItemMetaSchemaService(io.MPEndpointGetItemMetaSchemaInput{})
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusBadRequest, output)
		return
	}
	ctx.JSON(http.StatusOK, output)
}
func GetItemMetaHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	param := io.MPEndpointGetItemMetaInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	output := GetItemMetaService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, output)
}
func QueryItemMetasHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	param := io.MPEndpointQueryItemMetasInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	output := QueryItemMetasService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, output)
}
