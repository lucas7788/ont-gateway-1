package main

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/service"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	service.InitSellerImpl()
	seller.StartSellerServer()
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
