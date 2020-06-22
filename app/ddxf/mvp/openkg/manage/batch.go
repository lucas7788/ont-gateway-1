package main

import (
	"fmt"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/storage"
)

func main() {
	engine, err := storage.NewMySQL(config.Load().MySQLConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println(engine)
}
