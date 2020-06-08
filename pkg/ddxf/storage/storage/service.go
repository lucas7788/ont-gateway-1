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
	TestOntId string = "did:ont:Ad4pjz2bqep4RhQrUAzMuZJkBC3qJ1tZuT"
)

func UploadDataServiceHandle(c *gin.Context) {
	file, _, err := c.Request.FormFile("upload")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	output := UploadDataCore(file, TestOntId)
	c.JSON(output.Code, output)
}

func DownloadDataServiceHandle(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		instance.Logger().Error("UseTokenHandler:", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ResponseFailedOnto(http.StatusBadRequest, err))
		return
	}

	param := &io.StorageDownloadInput{}
	err = json.Unmarshal(data, &param)
	output := DowloadDataCore(param, TestOntId)
	c.JSON(output.Code, output)
}
