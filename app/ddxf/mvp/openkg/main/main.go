package main

import (
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
)

// MVP for openkg
func main() {
	common.ConsortiumAddr = "120.92.83.147:20336"
	server.StartServer()
}
