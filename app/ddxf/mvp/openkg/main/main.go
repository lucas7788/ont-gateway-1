package main

import (
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
)

// MVP for openkg
func main() {
	common.ConsortiumAddr = "113.31.112.154:20336"
	server.StartServer()
}
