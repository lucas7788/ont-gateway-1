package public

import (
	"github.com/gin-gonic/gin"
	g "github.com/zhiqiangxu/ont-gateway/pkg/rest/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
)

// NewApp for public app
func NewApp() *gin.Engine {
	r := g.New()

	addon := r.Group("/addon")

	attestation := addon.Group("/attestation")
	attestation.POST("/verify", middleware.AddonForward("/verify"))
	attestation.POST("/batchAdd", middleware.AddonForward("/batchAdd"))
	return r
}
