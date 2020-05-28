package main

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/registry/server"
	"os"
	"os/signal"
	"syscall"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
)

func main() {
	server.StartRegistryImplServer()
	waitToExit()
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			instance.Logger().Info("saga server received exit signal: %s." + sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
