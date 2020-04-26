package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// AkSkHello for test aksk
func AkSkHello(c *gin.Context) {
	var input io.AkSkHelloInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var output io.AkSkHelloOutput
	output.Msg = "ok"
	sendoutput(c, output.Code, output)
}
