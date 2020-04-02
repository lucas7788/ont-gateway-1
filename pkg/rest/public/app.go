package public

import (
	"github.com/gin-gonic/gin"
	g "github.com/zhiqiangxu/ont-gateway/pkg/rest/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
)

// NewApp for public app
func NewApp() *gin.Engine {
	r := g.New()

	biz := r.Group("/biz")

	attestation := biz.Group("/attestation")
	attestation.POST("/verify", middleware.AddonForward)
	attestation.POST("/batchAdd", middleware.AddonForward)
	return r
}
