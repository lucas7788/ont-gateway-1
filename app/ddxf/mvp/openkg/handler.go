package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
)

func publish(c *gin.Context) {
	var (
		input PublishInput
	)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	output, err := PublishService(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
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
		instance.DDXFSdk().DefDDXFKit().BuildBuyAndUseTokenTx()
		defer callback(output)
	}()
}
