package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
		GenerateOntIdService(input)
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
		PublishService(input)
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
