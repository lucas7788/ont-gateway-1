package metrics

import m "github.com/zhiqiangxu/util/metrics"

const (
	// RequestCount is name for request count
	RequestCount = "request_count"
	// RequestLatency is name for request latency
	RequestLatency = "request_latency"
)

// Register all metrics here
func Register() {
	m.RegisterCounter(RequestCount, []string{"method", "error"})
	m.RegisterHist(RequestLatency, []string{"method", "error"})
}
