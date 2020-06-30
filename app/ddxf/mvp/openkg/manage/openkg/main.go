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
	err = engine.SQL(`select b.id,b.url,b.description,b.format,a.author,a.maintainer,a.creator_user_id creatorID,b.resource_type,b.name, c.name author_name
	from 
			package a 
		left join resource b 
			on b.package_id=a.id 
		left join (select distinct name from tmp) c
			on a.author like concat('%',c.name,'%') or a.maintainer like concat('%',c.name,'%') 
	where 
		a.state='active' and 
		b.state='active'
	order by b.id 
		`).Find(&resources)
	if err != nil {
		panic(err)
	}

	resourceOwners := make(map[string]map[string]bool)
	for _, r := range resources {
		if resourceOwners[r.ID] == nil {
			resourceOwners[r.ID] = map[string]bool{}
		}
		resourceOwners[r.ID][r.Author] = true
		resourceOwners[r.ID][r.Maintainer] = true
	}

	for _, owners := range resourceOwners {
		for ownerName := range owners {
			if !ontid(ownerName) {
				panic(fmt.Sprintf("ontid fail for %s", ownerName))
			}
		}
	}

	resources = []resource{}
	err = engine.SQL(`select a.*, b.creator_user_id creatorID
	from 
			resource a
		left join package b 
			on b.id=a.package_id 
	where 
		a.state='active' and 
		b.state='active'
	order by b.id 
		`).Find(&resources)
	if err != nil {
		panic(err)
	}

	fmt.Println("resource count", len(resources))

	for i, r := range resources {

		fmt.Println("resource", r, "index", i)
		// r.CreatorID = "cd85f2c2-4fd1-44a8-82e3-10a7b63ed144"
		if !publish(r, resourceOwners[r.ID]) {
			panic(fmt.Sprintf("publish fail for %s, index:%d\n", r.ID, i))
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
	Author       string `xorm:"author"`
	Maintainer   string `xorm:"maintainer"`
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

func publish(r resource, owners map[string]bool) bool {

	dataMeta := map[string]interface{}{
		"id":            r.ID,
		"url":           r.URL,
		"description":   r.Description,
		"format":        r.Format,
		"author":        r.Author,
		"maintainer":    r.Maintainer,
		"creatorID":     r.CreatorID,
		"resource_type": r.ResourceType,
		"name":          r.Name,
	}

	var dataOwners []string
	for owner := range owners {
		dataOwners = append(dataOwners, owner)
	}

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
