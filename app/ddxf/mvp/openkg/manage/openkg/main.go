package main

import (
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/storage"
)

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
	fmt.Println(users)

	// for _, u := range users {
	// 	if !ontid(u.ID) {
	// 		panic(fmt.Sprintf("ontid fail for %s", u.ID))
	// 	}
	// 	fmt.Println("UserID", u.ID)
	// 	time.Sleep(time.Millisecond * 100)
	// }

	var resources []resource
	err = engine.SQL("select a.*, b.creator_user_id creatorID from resource a join package b on a.package_id=b.id").Find(&resources)
	if err != nil {
		panic(err)
	}
	fmt.Println(resources)

	for i, r := range resources {
		if i <= 690 {
			continue
		}
		fmt.Println("resource", r, "index", i)
		// r.CreatorID = "cd85f2c2-4fd1-44a8-82e3-10a7b63ed144"
		if !publish(r) {
			panic(fmt.Sprintf("publish fail for %s", r.ID))
		}
		// return
		time.Sleep(time.Millisecond * 100)
	}
}

type resource struct {
	ID           string `xorm:"id"`
	URL          string `xorm:"url"`
	Description  string `xorm:"description"`
	Format       string `xorm:"format"`
	CreatorID    string `xorm:"creatorID"`
	ResourceType string `xorm:"resource_type"`
	Name         string `xorm:"name"`
}

type user struct {
	ID string `xorm:"id"`
}

const domain = "http://openkg-dev.ontfs.io"

// const domain = "http://openkg-prod.ontfs.io"

// const domain = "http://192.168.0.228:10999"

func publish(r resource) bool {
	dataMeta := map[string]interface{}{
		"id":            r.ID,
		"url":           r.URL,
		"description":   r.Description,
		"format":        r.Format,
		"creatorID":     r.CreatorID,
		"resource_type": r.ResourceType,
		"name":          r.Name,
	}
	fmt.Println("r.CreatorID", r.CreatorID)
	input := server.PublishInput{
		ReqID:     uuid.NewV4().String(),
		OpenKGID:  r.ID,
		UserID:    r.CreatorID,
		OnChainId: uuid.NewV4().String(),
		Item:      dataMeta,
		Datas: []map[string]interface{}{
			dataMeta,
		},
	}

	bytes, _ := json.Marshal(input)
	code, _, _, err := forward.PostJSONRequest(domain+server.PublishURI, bytes, nil)
	return code == 200 && err == nil
}

func ontid(userID string) bool {
	input := server.GenerateOntIdInput{ReqID: uuid.NewV4().String(), UserId: userID}
	bytes, _ := json.Marshal(input)
	code, _, _, err := forward.PostJSONRequest(domain+server.GenerateOntIdByUserIdURI, bytes, nil)

	fmt.Println("code", code, "err", err)
	return code == 200 && err == nil
}
