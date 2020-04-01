package intra

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	g "github.com/zhiqiangxu/ont-gateway/pkg/rest/gin"
)

// NewApp for intra app
func NewApp() *gin.Engine {
	r := g.New()

	pprof.Register(r)
	// metrics内网展示
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
