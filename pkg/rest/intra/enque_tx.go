package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// EnqueTx for enqueue tx
func EnqueTx(c *gin.Context) {
	var input io.EnqueTxInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if input.Admin {
		c.JSON(http.StatusBadRequest, gin.H{"message": "admin denied"})
		return
	}

	input.App = c.GetInt(middleware.AppKey)

	output := service.Instance().EnqueTx(input)
	sendoutput(c, output.Code, output)
}
