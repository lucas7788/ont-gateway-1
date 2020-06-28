package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kataras/go-errors"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
	"go.uber.org/zap"
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

func GetDataIdByDataMetaHashHandler(c *gin.Context) {
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("[SaveDataMetaHandle] read param error: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	param := GetDataIdParam{}
	err = json.Unmarshal(bs, &param)
	if err != nil {
		instance.Logger().Error("[SaveDataMetaHandle] read param error: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	res, err := GetDataIdByDataMetaHashService(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func SaveDataMetaArrayHandler(c *gin.Context) {
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
	param := io.SellerSaveDataMetaArrayInput{}
	err = json.Unmarshal(bs, &param)
	if err != nil {
		instance.Logger().Error("[SaveDataMetaHandle] read param error: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	output := SaveDataMetaArrayService(param, ontId.(string))
	fmt.Println("SaveDataMetaArrayHandler:", output)
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
		instance.Logger().Error("[FreezeHandler] read param error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	param := DeleteParam{}
	err = json.Unmarshal(bs, &param)
	if err != nil {
		instance.Logger().Error("[FreezeHandler] unmarshal param error")
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	res := FreezeService(param, ontId.(string))
	c.JSON(res.Code, res)
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
		instance.Logger().Error("seller PublishMPItemMetaHandle:", zap.Error(output.Error()))
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, output.Error()))
		return
	}
	c.JSON(http.StatusOK, output)
}

func DeleteHandler(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("UseTokenHandler:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	param := DeleteInput{}
	err = json.Unmarshal(data, &param)
	if err != nil {
		instance.Logger().Error("UseTokenHandler:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	output := DeleteService(param)
	if output.Code == 0 {
		output.Code = http.StatusOK
	}
	c.JSON(output.Code, output)
}

func UpdateHandler(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("UseTokenHandler:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	param := UpdateInput{}
	err = json.Unmarshal(data, &param)
	if err != nil {
		instance.Logger().Error("UseTokenHandler:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	output := UpdateService(param)
	if output.Code == 0 {
		output.Code = http.StatusOK
	}
	c.JSON(output.Code, output)
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

func BuyAndUseDTokenHandler(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("BuyAndUseDTokenHandler:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	param := io.BuyerBuyAndUseDtokenInput{}
	err = json.Unmarshal(data, &param)
	if err != nil {
		instance.Logger().Error("BuyAndUseDTokenHandler:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	qrResp := BuyAndUseDTokenService(param)
	if qrResp.Error() != nil {
		instance.Logger().Error("BuyAndUseDTokenHandler:", zap.Error(qrResp.Error()))
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, qrResp.Error()))
		return
	}
	c.JSON(http.StatusOK, qrResp)
}
