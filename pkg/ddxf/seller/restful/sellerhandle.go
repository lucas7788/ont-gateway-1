package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/service"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
	//"go.uber.org/zap"
	"net/http"
)

func SaveDataMetaHandle(c *gin.Context) {
	ontId, ok := c.Get(middleware.TenantIDKey)
	if !ok {
		instance.Logger().Error("[SaveDataMetaHandle] read ontId error")
		//c.JSON(http.StatusBadRequest, ResponseFailedOnto(http.StatusBadRequest, err))
	}
	service.DefSellerImpl.SaveDataMeta()
}

func SaveTokenMetaHandle(c *gin.Context) {
	service.DefSellerImpl.SaveTokenMeta()
}

func PublishMPItemMetaHandle(c *gin.Context) {
	service.DefSellerImpl.PublishMPItemMeta()
}
