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

const (
	n = 1000

	collection = "openbase_research_graph_withID"

	batch = 1
)

var maintainers = map[string][]string{
	"openbase_research_graph_withID": []string{
		"陈卓", "吴杨", "邵鑫 (浙江大学药学院)", "杨海宏", "毕祯", "叶宏斌",
		"方尹", "杨帆", "陈华钧", "华为云语音语义创新Lab: 郑毅", "王鹏", // "卢栋才", "章涛", "袁晶",
		"华为云医疗智能体: 张雷", // "刘登辉", "徐迟", "乔楠",
	},
	"openbase_prevention_graph_withID": []string{
		"胡丹阳", "王萌", "李秋", "刘宇", "顾进广", "张志振", "胡闰秋", "胡闰秋", "张涛", "史淼", "郭文孜", "黄红蓝",
	},
	"openbase_medical_graph_withID": []string{
		"蔡嘉辉", "杜会芳",
	},
	"openbase_health_graph": []string{
		"许斌", "毛亦铭", "阎婧雅", "凤灵", "吴高晨", "仝美涵", "孙静怡", "李子明", "陈秋阳", "李凯曼", "郑晓飞", "刘邦长", "常德杰", "刘朝振", "刘红霞", "张航飞", "姜鹏", "闫广庆", "季科", "袁晓飞",
	},
	"openbase_goods_graph_withID": []string{
		"刘宇", "徐航", "向军毅", "顾进广",
	},
	"openbase_event_graph": []string{
		"刘作鹏", "王献敏", "彭茜", "戴振", "张作为", "王鲁威", "张呈阳", "刘杰", "唐彦",
	},
	"openbase_character_graph_withID": []string{
		"王智凤", "蔡嘉辉",
	},
	"openbase_wiki_graph": []string{
		"王昊奋", "漆桂林",
	},
}

var owners []string

func init() {
	owners = maintainers[collection]
	if len(owners) == 0 {
		panic(fmt.Sprintf("no maintainers for %s", collection))
	}
}

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

	cursor, err := instance.MongoOfficial().Collection(collection).Find(ctx, filter, options)
	if err != nil {
		panic(fmt.Sprintf("Find err:%v", err))
	}

	err = cursor.All(ctx, &resources)
	if err != nil {
		panic(fmt.Sprintf("cursor.All err:%v", err))
	}

	return
}

type user struct {
	Name string `xorm:"name"`
}

func main() {

	// var users []user
	// err = engine.SQL("select distinct name from tmp").Find(&users)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(users)

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

		ds := make([]map[string]interface{}, 0)
		for i, resource := range resources {

			ds = append(ds, resource)

			if len(ds) < batch {
				continue
			}

			fmt.Println("index", total-len(resources)+i)
			batchRegData(ds)
			// return

			// return
			time.Sleep(time.Millisecond * 1000)
		}

		if len(ds) > 0 {
			batchRegData(ds)
		}

	}

	fmt.Println("total", total)
}

const domain = "http://openkg-dev.ontfs.io"

// const domain = "http://openkg-prod.ontfs.io"

// const domain = "http://192.168.0.228:10999"

func ontid(userName string) bool {
	input := server.GenerateOntIdInput{ReqID: uuid.NewV4().String(), UserId: userName, Party: "openbase"}
	bytes, _ := json.Marshal(input)
	code, _, _, err := forward.PostJSONRequest(domain+server.GenerateOntIdByUserIdURI, bytes, nil)

	fmt.Println("code", code, "err", err)
	return code == 200 && err == nil
}

func batchRegData(ds []map[string]interface{}) {

	var partyDataIDs []string
	for _, d := range ds {
		partyDataIDs = append(partyDataIDs, d["@id"].(string))
	}
	var dataOwners [][]string
	for i := 0; i < len(ds); i++ {
		dataOwners = append(dataOwners, owners)
	}
	fmt.Println("partyDataIDs", partyDataIDs, "dataOwners", dataOwners)

	input := server.BatchRegDataInput{
		ReqID:        uuid.NewV4().String(),
		PartyDataIDs: partyDataIDs,
		Datas:        ds,
		DataOwners:   dataOwners,
		Party:        "openbase",
	}

	bytes, _ := json.Marshal(input)
	code, _, _, err := forward.PostJSONRequest(domain+server.BatchRegDataURI, bytes, nil)
	if !(code == 200 && err == nil) {
		panic(fmt.Sprintf("code:%v err:%v", code, err))
	}

}
