package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
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
	fileName := c.Param("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, "fileName can not empty")
		return
	}

	input := &io.StorageDownloadInput{
		FileName: fileName,
	}

	output := DowloadDataCore(input, TestOntId)
	c.JSON(output.Code, output)
}
