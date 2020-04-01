package main

import "github.com/zhiqiangxu/ont-gateway/pkg/rest"

func main() {

	server := rest.NewServer()
	server.Start()
}
