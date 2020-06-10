package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"net/http"
)

const (
	UploadDataUrl      = "/ddxf/storage/upload"
	DownloadDataPrefix = "/ddxf/storage/download/"
	DownloadDataUrl    = DownloadDataPrefix + ":fileName"
)

func StartStorageServer() {
	r := gin.Default()
	//r.Use(middleware.JWT)
	r.GET(config.Ping, func(context *gin.Context) {
		context.JSON(http.StatusOK, "SUCCESS")
	})
	r.POST(UploadDataUrl, UploadDataServiceHandle)
	r.GET(DownloadDataUrl, DownloadDataServiceHandle)
	go r.Run(":" + config.StorePort)
}
