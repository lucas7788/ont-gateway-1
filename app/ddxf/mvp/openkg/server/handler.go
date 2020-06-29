package server

import (
	"net/http"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kataras/go-errors"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
)

func GenerateOntIdByUserId(c *gin.Context) {
	var (
		input GenerateOntIdInput
	)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "SUCCESS")

	go func() {
		output := GenerateOntIdService(input)
		fmt.Println("GenerateOntIdByUserId output: ", output)
	}()
}

func Publish(c *gin.Context) {
	var (
		input PublishInput
	)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	go func() {
		output := PublishService(input)
		if output.Code != 0 {
			fmt.Println("openkg Publish:", output)
		}
	}()
	c.JSON(http.StatusOK, "SUCCESS")
}

func BuyAndUse(c *gin.Context) {
	var (
		input BuyAndUseInput
	)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	go func() {
		buyAndUseService(input)
	}()
	c.JSON(http.StatusOK, "SUCCESS")
}

func regDataHandler(c *gin.Context) {
	var (
		input RegDataInput
	)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	go func() {
		regDataService(input)
	}()
	c.JSON(http.StatusOK, "SUCCESS")
}
func deleteAttributesHandler(c *gin.Context) {
	var (
		input DeleteAttributesInput
	)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	go func() {
		output := deleteAttributesService(input)
		if output.Code != 0 {
			instance.Logger().Error("deleteAttributesService failed:", zap.Error(errors.New(output.Msg)))
			return
		}
	}()
	c.JSON(http.StatusOK, "SUCCESS")
}

func Delete(c *gin.Context) {
	var input DeleteInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	output := deleteService(input)
	if output.Code == 0 {
		output.Code = http.StatusOK
	}
	c.JSON(output.Code, output)
}
