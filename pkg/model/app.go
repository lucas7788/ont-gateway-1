package model

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zhiqiangxu/util"
)

// App for app
type App struct {
	ID               int    `bson:"id" json:"id"`
	Name             string `bson:"name" json:"name"`
	Ak               string `bson:"ak" json:"ak"`
	Sk               string `bson:"sk" json:"sk"`
	TxNotifyURL      string `bson:"tx_notify_url" json:"tx_notify_url"`
	PaymentNotifyURL string `bson:"payment_notify_url" json:"payment_notify_url"`
}

const (
	// GWAppName for ont-gateway
	GWAppName        = "ont-gw"
	maxContentLength = 1024 * 1024
)

type readerCloser struct {
	io.Reader
	io.Closer
}

// SignRequest generate aksk string for request
func (app *App) SignRequest(req *http.Request) (string, error) {
	h := hmac.New(sha1.New, util.Slice(app.Sk))

	var sb strings.Builder

	// 3 lines
	// 1st
	sb.WriteString(req.Method)
	sb.WriteString(" ")
	sb.WriteString(req.URL.Path)
	if req.URL.RawQuery != "" {
		sb.WriteString("?")
		sb.WriteString(req.URL.RawQuery)
	}
	// 2nd
	sb.WriteString("\nHost: ")
	sb.WriteString(req.Host)
	contentType := req.Header.Get("Content-Type")
	if contentType != "" {
		// 3rd
		sb.WriteString("\nContent-Type: ")
		sb.WriteString(contentType)
	}
	sb.WriteString("\n\n")

	if app.incBody(req, contentType) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "", err
		}

		sb.Write(body)

		req.Body = readerCloser{Reader: bytes.NewBuffer(body), Closer: req.Body}
	}

	str := sb.String()
	io.WriteString(h, str)
	encodeString := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return encodeString, nil
}

func (app *App) incBody(req *http.Request, ctType string) bool {
	typeOk := ctType != "" && ctType != "application/octet-stream"
	lengthOk := req.ContentLength > 0 && req.ContentLength < maxContentLength
	return typeOk && lengthOk && req.Body != nil
}
