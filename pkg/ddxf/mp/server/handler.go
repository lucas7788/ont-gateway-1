package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
	"io/ioutil"
)

func AddRegistryHandler(ctx *gin.Context) {
	handle(ctx, "AddRegistryHandler", io.MPAddRegistryInput{})
}

func RemoveRegistryHandler(ctx *gin.Context) {
	handle(ctx, "RemoveRegistryHandler", io.MPRemoveRegistryInput{})
}
func PublishItemMetaHandler(ctx *gin.Context) {
	handle(ctx, "PublishItemMetaHandler", io.MPEndpointPublishItemMetaInput{})
}

func GetAuditRuleHandler(ctx *gin.Context) {
	handle(ctx, "GetAuditRuleHandler", io.MPEndpointGetAuditRuleInput{})
}
func GetFeeHandler(ctx *gin.Context) {
	handle(ctx, "GetFeeHandler", io.MPEndpointGetFeeInput{})
}
func GetChallengePeriodHandler(ctx *gin.Context) {
	handle(ctx, "GetChallengePeriodHandler", io.MPEndpointGetChallengePeriodInput{})
}
func GetItemMetaSchemaHandler(ctx *gin.Context) {
	handle(ctx, "GetItemMetaSchemaHandler", io.MPEndpointGetItemMetaSchemaInput{})
}
func GetItemMetaHandler(ctx *gin.Context) {
	handle(ctx, "GetItemMetaHandler", io.MPEndpointGetItemMetaInput{})
}
func QueryItemMetasHandler(ctx *gin.Context) {
	handle(ctx, "QueryItemMetasHandler", io.MPEndpointQueryItemMetasInput{})
}

func handle(ctx *gin.Context, method string, param interface{}) {
	paramsBs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] read post param error:", zap.Error(err))
		return
	}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("[AddEndpointHandler] parse post param error:", zap.Error(err))
		return
	}
	switch method {
	case "AddRegistryHandler":
		output := AddRegistryService(param.(io.MPAddRegistryInput))
		ctx.JSON(output.Code, output)
	case "RemoveRegistryHandler":
		output := RemoveRegistryService(param.(io.MPRemoveRegistryInput))
		ctx.JSON(output.Code, output)
	case "PublishItemMetaHandler":
		output := PublishItemMetaService(param.(io.MPEndpointPublishItemMetaInput))
		ctx.JSON(output.Code, output)
	case "GetAuditRuleHandler":
		output := GetAuditRuleService(param.(io.MPEndpointGetAuditRuleInput))
		ctx.JSON(output.Code, output)
	case "GetFeeHandler":
		output := GetFeeService(param.(io.MPEndpointGetFeeInput))
		ctx.JSON(output.Code, output)
	case "GetChallengePeriodHandler":
		output := GetChallengePeriodService(param.(io.MPEndpointGetChallengePeriodInput))
		ctx.JSON(output.Code, output)
	case "GetItemMetaSchemaHandler":
		output := GetItemMetaSchemaService(param.(io.MPEndpointGetItemMetaSchemaInput))
		ctx.JSON(output.Code, output)
	case "GetItemMetaHandler":
		output := GetItemMetaService(param.(io.MPEndpointGetItemMetaInput))
		ctx.JSON(output.Code, output)
	case "QueryItemMetasHandler":
		output := QueryItemMetasService(param.(io.MPEndpointQueryItemMetasInput))
		ctx.JSON(output.Code, output)
	default:
		panic("not support service" + method)
	}
}
