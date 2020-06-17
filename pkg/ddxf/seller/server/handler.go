package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kataras/go-errors"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func SaveDataMetaHandler(c *gin.Context) {
	c.Set(middleware.TenantIDKey, "did:ont:AcVBV1zKGogf9Q54p1Ve78NSQVU5ZUUGkn")
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("[SaveDataMetaHandle] read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, errors.New("[SaveDataMetaHandle] read ontId error")))
		return
	}
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("[SaveDataMetaHandle] read param error: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	param := io.SellerSaveDataMetaInput{}
	err = json.Unmarshal(bs, &param)
	if err != nil {
		instance.Logger().Error("[SaveDataMetaHandle] read param error: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	output := SaveDataMetaService(param, ontId.(string))
	c.JSON(http.StatusOK, output)
}

func SaveTokenMetaHandler(c *gin.Context) {
	c.Set(middleware.TenantIDKey, "did:ont:AcVBV1zKGogf9Q54p1Ve78NSQVU5ZUUGkn")
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("[SaveTokenMetaHandler] read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, errors.New("[SaveTokenMetaHandler] read ontId error")))
		return
	}
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("[SaveTokenMetaHandler] read param error: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	param := io.SellerSaveTokenMetaInput{}
	err = json.Unmarshal(bs, &param)
	if err != nil {
		instance.Logger().Error("[SaveTokenMetaHandler] read param error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	output := SaveTokenMetaService(param, ontId.(string))
	if output.Code != 0 {
		instance.Logger().Error("[SaveTokenMetaHandler] read param error")
		c.JSON(output.Code, common.ResponseFailedOnto(common.PARA_ERROR, output.Error()))
		return
	}
	c.JSON(http.StatusOK, output)
}

func FreezeHandler(c *gin.Context) {
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("[SaveTokenMetaHandler] read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, errors.New("[SaveTokenMetaHandler] read ontId error")))
		return
	}
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("[PublishForOpenKgHandler] read param error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	param := FreezeParam{}
	err = json.Unmarshal(bs, &param)
	if err != nil {
		instance.Logger().Error("[PublishForOpenKgHandler] unmarshal param error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	res := FreezeService(param, ontId.(string))
	c.JSON(res.Code, res)
}

func PublishForOpenKgHandler(c *gin.Context) {
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("[SaveTokenMetaHandler] read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, errors.New("[SaveTokenMetaHandler] read ontId error")))
		return
	}
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("[PublishForOpenKgHandler] read param error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	param := PublishForOpenKgParam{}
	err = json.Unmarshal(bs, &param)
	if err != nil {
		instance.Logger().Error("[PublishForOpenKgHandler] unmarshal param error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	res := PublishForOpenKgService(param, ontId.(string))
	c.JSON(res.Code, res.Msg)
}

func PublishMPItemMetaHandler(c *gin.Context) {
	c.Set(middleware.TenantIDKey, "did:ont:AcVBV1zKGogf9Q54p1Ve78NSQVU5ZUUGkn")
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("PublishMPItemMetaHandle: read ontId error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, errors.New("PublishMPItemMetaHandle: read ontId error")))
		return
	}
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("[SaveDataMetaHandle] read param error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	param := io.MPEndpointPublishItemMetaInput{}
	err = json.Unmarshal(bs, &param)
	if err != nil {
		instance.Logger().Error("[SaveDataMetaHandle] read param error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	output := PublishMPItemMetaService(param, ontId.(string))
	if output.Error() != nil {
		instance.Logger().Error("PublishMPItemMetaHandle:", zap.Error(qrResp.Error()))
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, qrResp.Error()))
		return
	}
	c.JSON(http.StatusOK, output)
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
	qrResp := UseTokenService(param)
	if qrResp.Error() != nil {
		instance.Logger().Error("UseTokenHandler:", zap.Error(qrResp.Error()))
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, qrResp.Error()))
		return
	}
	c.JSON(http.StatusOK, qrResp)
}
