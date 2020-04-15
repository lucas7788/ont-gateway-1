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
func AddonForward(prefix string) func(*gin.Context) {
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

		resp, err := forward.Forward(c.Request, prefix, ad.SIP, "http")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			logger.Instance().Error("JSONRequest", zap.String("addonID", addonID), zap.String("tenantID", tenantID), zap.String("url", c.Request.URL.Path), zap.Error(err))
			return
		}

		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			logger.Instance().Error("JSONRequest", zap.String("url", c.Request.URL.Path), zap.Error(err))
			return
		}

		contentType := resp.Header.Get("Content-Type")
		c.DataFromReader(resp.StatusCode, int64(len(respBody)), contentType, bytes.NewReader(respBody), nil)
	}

}
