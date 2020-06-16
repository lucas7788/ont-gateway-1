package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

		// 2. save data metas and publish item
	}()

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
