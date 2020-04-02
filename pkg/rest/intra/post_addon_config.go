package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// PostAddonConfig for upsert a addon config
func PostAddonConfig(c *gin.Context) {

	var input io.PostAddonConfigInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	output := service.Instance().PostAddonConfig(input)
	sendoutput(c, output.Code, output)
}
