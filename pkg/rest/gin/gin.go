package gin

import (
	gin2 "github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/metrics"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
	m "github.com/zhiqiangxu/util/metrics"
)

// New returns a gin engine
func New() *gin2.Engine {
	g := gin2.New()

	requestLatencyMetric := m.GetHist(metrics.RequestLatency)
	requestCountMetric := m.GetCounter(metrics.RequestCount)

	mw := middleware.MetricLogger(requestLatencyMetric, requestCountMetric)
	g.Use(mw, gin2.Recovery())

	return g
}

func init() {
	metrics.Register()
}
