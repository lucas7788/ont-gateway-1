package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/middleware/cors"
)

const (
	AddEndpoint    = "/ddxf/buyer/BuyDtoken"
	RemoveEndpoint = "/ddxf/buyer/SaveTokenAndEndpoint"
	QueryEndpoint  = "/ddxf/buyer/UseToken"
)

func StartBuyerServer() {
	r := gin.Default()
	r.Use(cors.Cors())
	r.POST(AddEndpoint, BuyDtokenHandler)
	r.POST(QueryEndpoint, UseTokenHandler)
	Init()
	go r.Run(":" + "20331")
}