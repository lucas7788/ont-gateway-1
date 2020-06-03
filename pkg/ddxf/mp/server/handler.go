package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
	"io/ioutil"
)

func AddRegistryHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		return
	}
	param := io.MPAddRegistryInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		return
	}
	output := AddRegistryService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
	}
	ctx.JSON(output.Code, output)
}

func RemoveRegistryHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		return
	}
	param := io.MPRemoveRegistryInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		return
	}
	output := RemoveRegistryService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
	}
	ctx.JSON(output.Code, output)
}
func PublishItemMetaHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		return
	}
	param := io.MPEndpointPublishItemMetaInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		return
	}
	fmt.Println("param:", param)
	output := PublishItemMetaService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
	}
	ctx.JSON(output.Code, output)
}

func GetAuditRuleHandler(ctx *gin.Context) {
	output := GetAuditRuleService(io.MPEndpointGetAuditRuleInput{})
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
	}
	ctx.JSON(output.Code, output)
}
func GetFeeHandler(ctx *gin.Context) {
	output := GetFeeService(io.MPEndpointGetFeeInput{})
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
	}
	ctx.JSON(output.Code, output)
}
func GetChallengePeriodHandler(ctx *gin.Context) {
	output := GetChallengePeriodService(io.MPEndpointGetChallengePeriodInput{})
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
	}
	ctx.JSON(output.Code, output)
}
func GetItemMetaSchemaHandler(ctx *gin.Context) {
	output := GetItemMetaSchemaService(io.MPEndpointGetItemMetaSchemaInput{})
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
	}
	ctx.JSON(output.Code, output)
}
func GetItemMetaHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		return
	}
	param := io.MPEndpointGetItemMetaInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		return
	}
	output := GetItemMetaService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
	}
	ctx.JSON(output.Code, output)
}
func QueryItemMetasHandler(ctx *gin.Context) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		return
	}
	param := io.MPEndpointQueryItemMetasInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		return
	}
	output := QueryItemMetasService(param)
	if output.Code != 0 {
		instance.Logger().Error(output.Msg)
	}
	ctx.JSON(output.Code, output)
}
