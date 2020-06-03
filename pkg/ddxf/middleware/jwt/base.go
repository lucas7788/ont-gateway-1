package jwt

import (
	"encoding/base64"
	"encoding/json"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"time"
)

type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

type Payload struct {
	Aud     string  `json:"aud"`
	Iss     string  `json:"iss"`
	Exp     int     `json:"exp"`
	Iat     int     `json:"iat"`
	Jti     string  `json:"jti"`
	Content Content `json:"content"`
}

type Content struct {
	Type  string `json:"type"`
	OntId string `json:"ontId"`
	Role  string `json:"role"`
}

func GenerateJwt(userInfo string) (string, error) {
	h := Header{
		Typ: "JWT-X",
		Alg: "ES256",
		Kid: "",
	}
	c := Content{
		OntId: userInfo,
	}

	p := Payload{
		Exp:     int(time.Now().Unix() + int64(24*60*60)),
		Content: c,
	}

	hBs, err := json.Marshal(h)
	if err != nil {
		return "", err
	}
	pBs, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	raw := base64.RawURLEncoding.EncodeToString(hBs) + "." + base64.RawURLEncoding.EncodeToString(pBs)
	sig, err := config.DefDDXFConfig.OperatorAccount.Sign([]byte(raw))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(hBs) +
		"." + base64.RawURLEncoding.EncodeToString(pBs) +
		"." + base64.RawURLEncoding.EncodeToString(sig), nil
}
