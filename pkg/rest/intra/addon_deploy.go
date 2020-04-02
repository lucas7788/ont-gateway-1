package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// AddonDeploy for deploy a addon config
func AddonDeploy(c *gin.Context) {
	var input io.AddonDeployInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	output := service.Instance().AddonDeploy(input)
	sendoutput(c, output.Code, output)
}
