package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// CreatePaymentConfig for create payment config
func CreatePaymentConfig(c *gin.Context) {
	var input io.CreatePaymentConfigInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input.App = c.GetInt(middleware.AppKey)

	output := service.Instance().CreatePaymentConfig(input)
	sendoutput(c, output.Code, output)
}
