package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// DequeTx for dequeue tx
func DequeTx(c *gin.Context) {
	var input io.DequeTxInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if input.Admin {
		c.JSON(http.StatusBadRequest, gin.H{"message": "admin denied"})
		return
	}

	output := service.Instance().DequeTx(input)
	sendoutput(c, output.Code, output)
}
