package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"go.uber.org/zap"
)

// AddonForward for forward addon request
func AddonForward(forwardURI string) func(*gin.Context) {
	return func(c *gin.Context) {
		addonID := c.GetHeader("addonID")
		tenantID := c.GetHeader("tenantID")

		ad, exists, err := model.AddonDeploymentManager().Get(addonID, tenantID)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ad not found"})
			logger.Instance().Error("ad not found", zap.String("addonID", addonID), zap.String("tenantID", tenantID))
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			logger.Instance().Error("AddonDeploymentManager.Get", zap.String("addonID", addonID), zap.String("tenantID", tenantID), zap.Error(err))
			return
		}

		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		url := "http://" + ad.SIP + forwardURI
		code, contentType, respBytes, err := forward.JSONRequest(c.Request.Method, url, reqBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			logger.Instance().Error("JSONRequest", zap.String("url", url), zap.Error(err))
			return
		}

		c.DataFromReader(code, int64(len(respBytes)), contentType, bytes.NewReader(respBytes), nil)
	}

}
