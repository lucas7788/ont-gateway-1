package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// JSONLDAlign for align jsonld
func JSONLDAlign(c *gin.Context) {
	var input io.JSONLDAlignInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	output := service.Instance().JSONLDAlign(input)
	sendoutput(c, output.Code, output)
}
