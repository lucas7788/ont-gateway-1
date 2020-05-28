package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/registry/middleware/cors"
)

const (
	AddEndpoint    = "/ddxf/registry/addendpoint"
	RemoveEndpoint = "/ddxf/registry/removeendpoint"
	QueryEndpoint  = "/ddxf/registry/queryendpoint"
)

func StartRegistryImplServer() {
	r := gin.Default()
	r.Use(cors.Cors())
	r.POST(AddEndpoint, AddEndpointHandler)
	r.POST(RemoveEndpoint, RemoveEndpointHandler)
	r.GET(QueryEndpoint, QueryEndpointHandler)
	Init()
	go r.Run(":" + "20331")
}
