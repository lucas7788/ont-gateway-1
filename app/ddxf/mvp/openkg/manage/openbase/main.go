package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const n = 1000

func getResource(page int) (resources []map[string]interface{}) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{}
	options := options.Find()
	// sort by _id desc
	options.SetSort(bson.D{bson.E{Key: "_id", Value: -1}})
	options.SetSkip(int64(page * n))
	options.SetLimit(n)

	cursor, err := instance.MongoOfficial().Collection("openbase_research_graph_withID").Find(ctx, filter, options)
	if err != nil {
		panic(fmt.Sprintf("Find err:%v", err))
	}

	err = cursor.All(ctx, &resources)
	if err != nil {
		panic(fmt.Sprintf("cursor.All err:%v", err))
	}

	return
}

func main() {

	total := 0
	page := 0
	for {

		resources := getResource(page)
		if len(resources) == 0 {
			break
		}

		total += len(resources)
		fmt.Println("page", page, "total", total)
		page++

		for i, resource := range resources {
			fmt.Println("index", total-len(resources)+i)
			regData(resource)

			// return
			time.Sleep(time.Millisecond * 100)
		}

	}

	fmt.Println("total", total)
}

const domain = "http://openkg-dev.ontfs.io"

// const domain = "http://openkg-prod.ontfs.io"

// const domain = "http://192.168.0.228:10999"

func regData(d map[string]interface{}) {
	input := server.RegDataInput{
		ReqID:       uuid.NewV4().String(),
		PartyDataID: d["@id"].(string),
		Data:        d,
		Controllers: []string{},
		Party:       "openbase",
	}

	bytes, _ := json.Marshal(input)
	code, _, _, err := forward.PostJSONRequest(domain+server.PublishURI, bytes, nil)
	if !(code == 200 && err == nil) {
		panic(fmt.Sprintf("code:%v err:%v", code, err))
	}
}

// func publish(r resource) bool {
// 	dataMeta := map[string]interface{}{
// 		"id":            r.ID,
// 		"url":           r.URL,
// 		"description":   r.Description,
// 		"format":        r.Format,
// 		"creatorID":     r.CreatorID,
// 		"resource_type": r.ResourceType,
// 		"name":          r.Name,
// 	}
// 	fmt.Println("r.CreatorID", r.CreatorID)
// 	input := server.PublishInput{
// 		ReqID:     uuid.NewV4().String(),
// 		OpenKGID:  r.ID,
// 		UserID:    r.CreatorID,
// 		OnChainId: uuid.NewV4().String(),
// 		Item:      dataMeta,
// 		Datas: []map[string]interface{}{
// 			dataMeta,
// 		},
// 	}

// 	bytes, _ := json.Marshal(input)
// 	code, _, _, err := forward.PostJSONRequest(domain+server.PublishURI, bytes, nil)
// 	return code == 200 && err == nil
// }
