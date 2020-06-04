package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func GetQrCodeDataByQrCodeIdHandler(c *gin.Context) {
	qrCodeId := c.Param("qrCodeId")
	code, err := GetQrCodeByQrCodeIdService(qrCodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		c.JSON(http.StatusOK, common.ResponseSuccess(code))
	}
}

func GrCodeCallbackSendTxHandler(c *gin.Context) {
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("GrCodeCallbackSendTxHandler: read post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	param := qrCode.QrCodeCallBackParam{}
	err = json.Unmarshal(paramsBs, &param)
	if err != nil {
		instance.Logger().Error("GrCodeCallbackSendTxHandler: parse post param error:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusInternalServerError, err))
		return
	}
	err = QrCodeCallBackService(param)
	if err != nil {
		instance.Logger().Error("GrCodeCallbackSendTxHandler: :", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.ResponseFailedOnto(http.StatusInternalServerError, err))
	} else {
		c.JSON(http.StatusOK, common.ResponseSuccess("SUCCESS"))
	}
}
