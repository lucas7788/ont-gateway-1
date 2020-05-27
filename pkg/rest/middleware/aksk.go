package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

const (
	// AppKey for app
	AppKey = "app"
)

// AkSk middleware
func AkSk(c *gin.Context) {

	authString := c.GetHeader("Authorization")
	if authString == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Authorization not found"})
		c.Abort()
		return
	}
	appStr := c.GetHeader("App")
	if appStr == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "App not found"})
		c.Abort()
		return
	}
	appInt64, err := strconv.ParseInt(appStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "App invalid"})
		c.Abort()
		return
	}
	appInt := int(appInt64)

	app := model.AppManager().GetByID(appInt)
	if app == nil {
		c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("App not found:%v", appInt)})
		c.Abort()
		return
	}

	// authString is like: mt ak:signature
	// extract ak
	spaceIdx := strings.IndexByte(authString, ' ')
	if spaceIdx == -1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Authorization not valid"})
		c.Abort()
		return
	}
	colonIdx := strings.IndexByte(authString[spaceIdx:], ':')
	if colonIdx == -1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Authorization not valid"})
		c.Abort()
		return
	}

	skIdx := spaceIdx + 1 + colonIdx
	ak := authString[spaceIdx+1 : skIdx-1]

	if app.Ak != ak {
		c.JSON(http.StatusForbidden, gin.H{"message": "App and Ak not match"})
		c.Abort()
		return
	}
	c.Set(AppKey, app.ID)

	// verify
	sign, err := app.SignRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "SignRequest err"})
		c.Abort()
		return
	}

	if sign != authString[skIdx:] {
		instance.Logger().Error("aksk not match", zap.String("expect", sign), zap.String("got", authString[skIdx:]))
		c.JSON(http.StatusForbidden, gin.H{"message": "aksk sign not match"})
		c.Abort()
		return
	}

	c.Next()
}
