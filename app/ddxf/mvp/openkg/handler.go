package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/key_manager"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
)

func publish(c *gin.Context) {
	var (
		input  PublishInput
		output PublishOutput
	)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	go func() {
		output.ReqID = input.ReqID
		defer callback(output)

		// 抽取openKGID
		openKGID := input.OpenKGID

		if input.Delete {

			// 下架
		} else {
		}

		// 1. 抽取data meta
		dataMetas := input.Datas
		plainSeed := []byte(defPlainSeed + input.UserID)
		pri,_ := key_manager.GetSerializedKeyPair(plainSeed)
		seller,err := ontology_go_sdk.NewAccountFromPrivateKey(pri,signature.SHA256withECDSA)
		if err != nil {
			return
		}
		resourceId := common.GenerateUUId(config.UUID_RESOURCE_ID)
		// 2. save data metas and publish item
		// send request to seller
		instance.DDXFSdk().SetPayer(payer)

		tx, err := instance.DDXFSdk().DefDDXFKit().BuildPublishTx([]byte(resourceId),)
		instance.DDXFSdk().SignTx(tx, seller)
	}()

}

func freeze(c *gin.Context) {
	instance.DDXFSdk().DefDDXFKit().Freeze()
}

func buyAndUse(c *gin.Context) {
	var (
		input  BuyAndUseInput
		output BuyAndUseOutput
	)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	go func() {
		output.ReqID = input.ReqID
		defer callback(output)

	}()
}
