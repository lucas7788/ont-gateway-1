package main

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/example/seller_buyer"
)

var (
	dataMetaHash  = "f1dfe4c60f9f8e4942559ee14c549ce63abfced6a1be08519761744f2429ac35"
	tokenMetaHash = "e2a740fa12bd94f0e242688e29f6d803f7671eb1f81bcfbdc1c3e213878e7dd4"
	resourceId    = ""
	dataId        = "dataid_1d3294cd-b30e-460a-a358-d02fb3d01699"
)

func main() {
	if true {
		err := seller_buyer.SaveDataMeta()
		if err != nil {
			panic(err)
		}
		return
	}
	if false {
		err := seller_buyer.SaveTokenMeta(dataMetaHash)
		if err != nil {
			panic(err)
		}
		return
	}
	if true {
		err := seller_buyer.PublishMeta1(tokenMetaHash, dataMetaHash)
		if err != nil {
			panic(err)
		}
		return
	}
	if true {
		err := seller_buyer.PublishMeta()
		if err != nil {
			panic(err)
		}
		return
	}
	if true {
		err := seller_buyer.BuyDtoken(resourceId)
		if err != nil {
			panic(err)
		}
		return
	}
	if true {
		err := seller_buyer.UseToken(resourceId, tokenMetaHash, dataId)
		if err != nil {
			panic(err)
		}
		return
	}
}
