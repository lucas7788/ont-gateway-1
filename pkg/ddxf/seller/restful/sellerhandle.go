package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/service"
)

func SaveDataMetaHandle(c *gin.Context) {
	service.DefSellerImpl.SaveDataMeta()
}

func SaveTokenMetaHandle(c *gin.Context) {
	service.DefSellerImpl.SaveTokenMeta()
}

func PublishMPItemMetaHandle(c *gin.Context) {
	service.DefSellerImpl.PublishMPItemMeta()
}
