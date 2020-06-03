package restful

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/service"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func SaveDataMetaHandle(c *gin.Context) {
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("[SaveDataMetaHandle] read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, nil))
		return
	}
	param := io.SellerSaveDataMetaInput{}
	output := service.DefSellerImpl.SaveDataMeta(param, ontId.(string))
	c.JSON(output.Code, output)
}

func SaveTokenMetaHandle(c *gin.Context) {
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("[SaveDataMetaHandle] read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, nil))
		return
	}
	param := io.SellerSaveTokenMetaInput{}
	output := service.DefSellerImpl.SaveTokenMeta(param, ontId.(string))
	c.JSON(output.Code, output)
}

func PublishMPItemMetaHandle(c *gin.Context) {
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("PublishMPItemMetaHandle: read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, nil))
		return
	}
	param := io.MPEndpointPublishItemMetaInput{}

	qrResp := service.DefSellerImpl.PublishMPItemMeta(param, ontId.(string))
	if qrResp.Error() != nil {
		instance.Logger().Error("PublishMPItemMetaHandle:", zap.Error(qrResp.Error()))
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, qrResp.Error()))
		return
	}
	c.JSON(http.StatusOK, qrResp)
}

func UseTokenHandler(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("UseTokenHandler:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	param := io.SellerTokenLookupEndpointUseTokenInput{}
	err = json.Unmarshal(data, &param)
	if err != nil {
		instance.Logger().Error("UseTokenHandler:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	qrResp := service.UseTokenService(param)
	if qrResp.Error() != nil {
		instance.Logger().Error("UseTokenHandler:", zap.Error(qrResp.Error()))
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, qrResp.Error()))
		return
	}
	c.JSON(http.StatusOK, qrResp)
}

func PublishMetaHandler(c *gin.Context) {
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("PublishMetaHandler: read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, nil))
		return
	}
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("PublishMetaHandler: read post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	param := io.SellerPublishMPItemMetaInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("PublishMetaHandler: parse post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	res, err := service.PublishMetaService(param, ontId.(string))
	if err != nil {
		instance.Logger().Error("PublishMetaHandler: parse post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetQrCodeDataByQrCodeId(c *gin.Context) {
	qrCodeId := c.Param("qrCodeId")
	code, err := service.GetQrCodeByQrCodeIdService(qrCodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		c.JSON(http.StatusOK, common.ResponseSuccess(code))
	}
}

func GrCodeCallbackSendTx(c *gin.Context) {
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("GrCodeCallbackSendTx: read post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	param := qrCode.QrCodeCallBackParam{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("GrCodeCallbackSendTx: parse post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	err = service.QrCodeCallBackService(param)
	if err != nil {
		instance.Logger().Error("GrCodeCallbackSendTx: :", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		c.JSON(http.StatusOK, common.ResponseSuccess("SUCCESS"))
	}
}
