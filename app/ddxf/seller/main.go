package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
)

func main() {
	common.ConsortiumAddr = "113.31.112.154:20336"
	server.StartSellerServer()
	waitToExit()
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			instance.Logger().Error("ddxf seller server received exit signal: ." + sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
