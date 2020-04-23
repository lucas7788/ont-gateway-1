package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// CreatePaymentOrder for create payment order
func CreatePaymentOrder(c *gin.Context) {
	var input io.CreatePaymentOrderInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	output := service.Instance().CreatePaymentOrder(input)
	sendoutput(c, output.Code, output)
}
