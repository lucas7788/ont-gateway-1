package main

import (
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/example/seller_buyer"
)

var (
	dataMetaHash  = "e2f09f0575858f3e09b5362b72a52e24b20e52828865003429b2f17539a686f6"
	tokenMetaHash = "e2a740fa12bd94f0e242688e29f6d803f7671eb1f81bcfbdc1c3e213878e7dd4"
	resourceId    = "resourceid_87b45ad3-0199-4c76-87e2-837b86e27bd2"
	dataId        = "dataid_2e214f27-e599-4fc3-9b3f-c62f2fb12464"
	qrCodeId      = "seller_publish886fc0d4-d5b6-4bd2-8c46-28ea395cc90e"
)

func main() {
	if false {
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
	if false {
		err := seller_buyer.PublishMeta1(tokenMetaHash, dataMetaHash)
		if err != nil {
			panic(err)
		}
		return
	}
	if false {
		err := seller_buyer.PublishMeta(qrCodeId)
		if err != nil {
			panic(err)
		}
		return
	}
	pwd := []byte("123456")
	wallet, _ := ontology_go_sdk.OpenWallet("/Users/sss/gopath/src/github.com/zhiqiangxu/ont-gateway/pkg/ddxf/example/wallet.dat")
	buyer, _ := wallet.GetAccountByAddress("AHhXa11suUgVLX1ZDFErqBd3gskKqLfa5N", pwd)
	if false {
		err := seller_buyer.BuyDtoken(buyer, resourceId)
		if err != nil {
			fmt.Println("error:", err)
		}
		return
	}
	if true {
		err := seller_buyer.UseToken(buyer, resourceId, tokenMetaHash, dataId)
		if err != nil {
			panic(err)
		}
		return
	}
}
