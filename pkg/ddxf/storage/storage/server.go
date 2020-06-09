package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
)

const (
	UploadDataUrl      = "/ddxf/storage/upload"
	DownloadDataPrefix = "/ddxf/storage/download/"
	DownloadDataUrl    = DownloadDataPrefix + ":fileName"
)

func StartStorageServer() {
	r := gin.Default()
	//r.Use(middleware.JWT)
	r.POST(UploadDataUrl, UploadDataServiceHandle)
	r.GET(DownloadDataUrl, DownloadDataServiceHandle)
	go r.Run(":" + config.StorePort)
}
