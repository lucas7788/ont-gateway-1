package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// DeletePaymentConfig for delete payment config
func DeletePaymentConfig(c *gin.Context) {
	var input io.DeletePaymentConfigInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	output := service.Instance().DeletePaymentConfig(input)
	sendoutput(c, output.Code, output)
}
