package server

import (
	"github.com/gin-gonic/gin"
)

const (
	AddEndpoint = "/ddxf/registry/addendpoint"
	RemoveEndpoint = "/ddxf/registry/removeendpoint"
	QueryEndpoint = "/ddxf/registry/queryendpoint"
)

func StartRegistryImplServer() {
	r := gin.Default()
	r.POST(AddEndpoint, AddEndpointHandler)
	r.POST(RemoveEndpoint, RemoveEndpointHandler)
	r.GET(QueryEndpoint, QueryEndpointHandler)
	go r.Run(":" + "20331")
}