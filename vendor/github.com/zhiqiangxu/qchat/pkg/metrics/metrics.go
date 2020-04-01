package metrics

import "github.com/zhiqiangxu/util/metrics"

const (
	// RequestCount is name for request count
	RequestCount = "request_count"
	// RequestLatency is name for request latency
	RequestLatency = "request_latency"
	// OnlineCount for online count
	OnlineCount = "online_count"
	// KickCount for kick count
	KickCount = "kick_count"
	// NotifyPendingCount for notify pending count
	NotifyPendingCount = "notify_pending_count"
)

// Register all metrics here
func Register() {
	metrics.RegisterCounter(RequestCount, []string{"method", "error"})
	metrics.RegisterHist(RequestLatency, []string{"method", "error"})
	metrics.RegisterGauge(OnlineCount, []string{"app", "kind"})
	metrics.RegisterGauge(NotifyPendingCount, []string{})
	metrics.RegisterCounter(KickCount, []string{"app"})
}
