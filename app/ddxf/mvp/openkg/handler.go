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
	ui, err := GenerateOntIdService(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ui)
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
