package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// UpsertAddonConfig for upsert a addon config
func UpsertAddonConfig(c *gin.Context) {

	var input io.UpsertAddonConfigInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	output := service.Instance().UpsertAddonConfig(input)
	sendoutput(c, output.Code, output)
}
