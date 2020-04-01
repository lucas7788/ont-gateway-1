package public

import (
	"github.com/gin-gonic/gin"
	g "github.com/zhiqiangxu/ont-gateway/pkg/rest/gin"
)

// NewApp for public app
func NewApp() *gin.Engine {
	r := g.New()

	return r
}
