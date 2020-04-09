package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
)

const (
	// MPKey for marketplace
	MPKey = "mp"
	// TenantIDKey for tenant id
	TenantIDKey = "tid"
	// AuthorizationKey for Authorization
	AuthorizationKey = "Authorization"
)

// JWT middleware
func JWT(c *gin.Context) {

	authString := c.GetHeader(AuthorizationKey)
	if authString == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Authorization not found"})
		c.Abort()
		return
	}

	jwt := instance.JWT()
	ok, claims := jwt.Validate(authString)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"message": "jwt validation failed"})
		c.Abort()
		return
	}

	mp, ok := claims[MPKey]
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"message": "mp not in claims"})
		c.Abort()
		return
	}
	tid, ok := claims[TenantIDKey]
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"message": "tid not in claims"})
		c.Abort()
		return
	}

	c.Set(MPKey, mp)
	c.Set(TenantIDKey, tid)

	c.Next()
}
