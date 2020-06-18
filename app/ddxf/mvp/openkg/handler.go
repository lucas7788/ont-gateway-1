package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateOntIdByUserId(c *gin.Context) {
	var (
		input GenerateOntIdInput
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

func publish(c *gin.Context) {
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

func buyAndUse(c *gin.Context) {
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
