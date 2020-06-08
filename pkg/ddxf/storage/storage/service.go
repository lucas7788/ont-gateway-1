package storage

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

const (
	TestOntId = "did:ont:Ad4pjz2bqep4RhQrUAzMuZJkBC3qJ1tZuT"
	UploadKey = "upload"
)

func UploadDataServiceHandle(c *gin.Context) {
	file, _, err := c.Request.FormFile(UploadKey)
	if err != nil {
		instance.Logger().Error("[UploadDataServiceHandle] error: ", zap.Error(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	output := UploadDataCore(file, TestOntId)
	c.JSON(output.Code, output)
}

func DownloadDataServiceHandle(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("DownloadDataServiceHandle:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}

	param := &io.StorageDownloadInput{}
	err = json.Unmarshal(data, &param)
	if err != nil {
		instance.Logger().Error("DownloadDataServiceHandle:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}
	output := DowloadDataCore(param, TestOntId)
	c.JSON(output.Code, output)
}
