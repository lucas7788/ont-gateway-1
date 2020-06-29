package main

import (
	"fmt"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/storage"
)

type user struct {
	ID string `xorm:"id"`
}

func main() {
	engine, err := storage.NewMySQL(config.Load().MySQLConfig)
	if err != nil {
		panic(err)
	}

	var users []user
	err = engine.SQL("select * from user").Find(&users)
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		fmt.Println(user)
	}
}
