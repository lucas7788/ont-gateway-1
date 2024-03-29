package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/middleware/cors"
	"net/http"
)

const (
	AddEndpoint    = "/ddxf/registry/addendpoint"
	RemoveEndpoint = "/ddxf/registry/removeendpoint"
	QueryEndpoint  = "/ddxf/registry/queryendpoint"
)

func StartRegistryImplServer() {
	r := gin.Default()
	r.Use(cors.Cors())
	r.GET(config.Ping, func(context *gin.Context) {
		context.JSON(http.StatusOK, "SUCCESS")
	})
	r.POST(AddEndpoint, AddEndpointHandler)
	r.POST(RemoveEndpoint, RemoveEndpointHandler)
	r.GET(QueryEndpoint, QueryEndpointHandler)
	Init()
	go r.Run(":" + config.RegistryPort)
}
