package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
	"go.uber.org/zap"
)

func GetQrCodeDataByQrCodeIdHandlerOnto(c *gin.Context) {
	qrCodeId := c.Param("qrCodeId")
	code, err := GetQrCodeByQrCodeIdService(qrCodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		c.JSON(http.StatusOK, common.ResponseSuccess(code))
	}
}

func GrCodeCallbackSendTxHandlerOnto(c *gin.Context) {
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("GrCodeCallbackSendTxHandlerOnto: read post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	param := qrCode.QrCodeCallBackParam{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("GrCodeCallbackSendTxHandlerOnto: parse post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	err = QrCodeCallBackService(param)
	if err != nil {
		instance.Logger().Error("GrCodeCallbackSendTxHandlerOnto: :", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		c.JSON(http.StatusOK, common.ResponseSuccess("SUCCESS"))
	}
}

func PublishMetaHandlerOnto(c *gin.Context) {
	c.Set(middleware.TenantIDKey, "did:ont:AcVBV1zKGogf9Q54p1Ve78NSQVU5ZUUGkn")
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("PublishMetaHandlerOnto: read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, errors.New("PublishMPItemMetaHandle: read ontId error")))
		return
	}
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("PublishMetaHandlerOnto: read post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	param := io.SellerPublishMPItemMetaInput{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("PublishMetaHandlerOnto: parse post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	res, err := PublishMetaService(param, ontId.(string))
	if err != nil {
		instance.Logger().Error("PublishMetaHandlerOnto: parse post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, res)
}
