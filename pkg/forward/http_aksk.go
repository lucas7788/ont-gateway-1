package forward

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

func getAuthorization(app *model.App, method, host, uri string, contentType string, body []byte) string {

	// 第一行Method Path
	firstLine := method + " " + uri + "\n"
	secondLine := "Host: " + host + "\n"
	thirdLine := "Content-Type: " + contentType + "\n\n"

	str := firstLine + secondLine + thirdLine + string(body)

	//hmac ,use sha1
	key := []byte(app.Sk)
	h := hmac.New(sha1.New, key)

	io.WriteString(h, str)

	encodeString := base64.URLEncoding.EncodeToString(h.Sum(nil))

	code := "mt " + app.Ak + ":" + encodeString
	return code
}

var (
	contentType = "application/json"
	postMethod  = "POST"
)

// PostAkSkRequestByName for post aksk request by app name
func PostAkSkRequestByName(appName, host, uri string, data interface{}) (code int, respContentType string, respBody []byte, err error) {
	app := model.AppManager().GetByName(appName)
	if app == nil {
		err = fmt.Errorf("app not exists:%s", appName)
		return
	}

	return PostAkSkRequest(app, host, uri, data)
}

// PostAkSkRequest for send a Post aksk request
func PostAkSkRequest(app *model.App, host, uri string, data interface{}) (code int, respContentType string, respBody []byte, err error) {
	return AkSkRequest(app, postMethod, host, uri, data)
}

// AkSkRequest send an aksk request and returns response
func AkSkRequest(app *model.App, method string, host, uri string, data interface{}) (code int, respContentType string, respBody []byte, err error) {

	body, err := json.Marshal(data)

	if err != nil {
		return
	}

	url := "http://" + host + uri
	req, err := http.NewRequest(method, url, strings.NewReader(string(body)))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", getAuthorization(app, method, host, uri, contentType, body))

	return httpRequest(req)

}
