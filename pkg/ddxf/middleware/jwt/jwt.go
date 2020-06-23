package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/core/signature"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"net/http"
	"strings"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var load *Payload
		header := map[string][]string(c.Request.Header)
		token := header["Authorization"]
		if token == nil || token[0] == "" {
			err = fmt.Errorf("token is nil")
		} else {
			load, err = validateToken(token[0], false)
		}
		if err != nil {
			instance.Logger().Error("token error: " + err.Error())
			c.JSON(http.StatusUnauthorized, common.ResponseFailed(common.VERIFY_TOKEN_ERROR, err))
			c.Abort()
			return
		}
		c.Set(config.Key_OntId, load.Content.OntId)
		c.Set(config.JWTAud, load.Aud)
		c.Set(config.JWTAdmin, false)
		c.Next()
	}
}

func validateToken(token string, admin bool) (*Payload, error) {
	//header.payloadBs.sig
	arr := strings.Split(token, ".")
	if len(arr) != 3 {
		return nil, fmt.Errorf("wrong token: %s", token)
	}
	sig, err := base64.RawURLEncoding.DecodeString(arr[2])
	if err != nil {
		return nil, err
	}

	data := arr[0] + "." + arr[1]
	err = signature.Verify(config.DefDDXFConfig().OperatorAccount.GetPublicKey(), []byte(data), sig)
	if err != nil {
		return nil, err
	}
	payloadBs, err := base64.RawURLEncoding.DecodeString(arr[1])
	if err != nil {
		return nil, err
	}
	pl := &Payload{}
	err = json.Unmarshal(payloadBs, pl)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	if pl.Exp < int(now) {
		return nil, fmt.Errorf("jwt token expired")
	}
	return pl, nil
}
