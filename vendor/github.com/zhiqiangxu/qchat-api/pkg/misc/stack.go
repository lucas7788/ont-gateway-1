package misc

import (
	"runtime"
	"github.com/zhiqiangxu/qrpc"
)

// StackTrace for stack trace
func StackTrace() string {
	// Reserve 10K buffer
	buf := make([]byte, 10240)
	buf = buf[:runtime.Stack(buf, false)]
	return qrpc.String(buf)
}
