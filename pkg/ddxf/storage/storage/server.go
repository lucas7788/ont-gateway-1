package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
)

const (
	UploadDataUrl   = "/ddxf/storage/upload"
	DownloadDataUrl = "/ddxf/storage/download"
)

func StartStorageServer() {
	r := gin.Default()
	//r.Use(middleware.JWT)
	r.POST(UploadDataUrl, UploadDataServiceHandle)
	r.POST(DownloadDataUrl, DownloadDataServiceHandle)
	go r.Run(":" + config.StorePort)
}
