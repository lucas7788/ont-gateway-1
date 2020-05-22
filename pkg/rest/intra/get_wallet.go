package intra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// GetWallet for get wallet info
func GetWallet(c *gin.Context) {
	// clientIP := c.ClientIP()
	// whiteList := []string{
	// 	"103.61.36.254",
	// 	// 公司的IP
	// 	"54.255.218.249",
	// 	// 生产build
	// 	"52.198.33.251",
	// }
	// whiteListMap := make(map[string]struct{})
	// for _, ip := range whiteList {
	// 	whiteListMap[ip] = struct{}{}
	// }
	// if _, ok := whiteListMap[clientIP]; !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("access denied for %s", clientIP)})
	// 	return
	// }

	var input io.GetWalletInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	output := service.Instance().GetWallet(input)
	sendoutput(c, output.Code, output)
}
