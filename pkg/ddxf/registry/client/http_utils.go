package client

import (
	"net/http"
	"time"
	"strings"
	"net/url"
	"fmt"
	"encoding/json"
	"io"
	"bytes"
	"io/ioutil"
)

type HttpClient struct {
	addr       string
	httpClient *http.Client
}


func NewHttpClient(addr string) *HttpClient {
	return &HttpClient{
		addr:addr,
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   5,
				DisableKeepAlives:     false, //enable keepalive
				IdleConnTimeout:       time.Second * 300,
				ResponseHeaderTimeout: time.Second * 300,
			},
			Timeout: time.Second * 300, //timeout for http response
		},
	}
}

func (this *HttpClient) getRequestUrl(reqPath string, values ...*url.Values)(string, error) {
	var addr string
	if !strings.HasPrefix(this.addr, "http") {
		addr = "http://" + addr
	}
	reqUrl, err := new(url.URL).Parse(addr)
	if err != nil {
		return "", fmt.Errorf("Parse address:%s error:%s", addr, err)
	}
	reqUrl.Path = reqPath
	if len(values) > 0 && values[0] != nil {
		reqUrl.RawQuery = values[0].Encode()
	}
	return reqUrl.String(), nil
}

func (this *HttpClient) SendPostRequest(reqParam interface{}, reqPath string, values ...*url.Values)([]byte, error) {
	reqUrl, err := this.getRequestUrl(reqPath, values...)
	if err != nil {
		return nil, err
	}
	reqData, err := json.Marshal(reqParam)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal error:%s", err)
	}
	resp, err := this.httpClient.Post(reqUrl, "application/json", bytes.NewReader(reqData))
	if err != nil {
		return nil, fmt.Errorf("send http post request error:%s", err)
	}
	defer resp.Body.Close()
	return this.dealRestResponse(resp.Body)
}

func (this *HttpClient) dealRestResponse(body io.Reader) ([]byte, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read http body error:%s", err)
	}
	restRsp := &RestfulResp{}
	err = json.Unmarshal(data, restRsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal RestfulResp:%s error:%s", body, err)
	}
	if restRsp.Error != 0 {
		return nil, fmt.Errorf("sendRestRequest error code:%d desc:%s result:%s", restRsp.Error, restRsp.Desc, restRsp.Result)
	}
	return restRsp.Result, nil
}

func (this *HttpClient) SendGetRequest(reqPath string, values ...*url.Values) ([]byte, error) {
	reqUrl, err := this.getRequestUrl(reqPath, values...)
	if err != nil {
		return nil, err
	}
	resp, err := this.httpClient.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("send http get request error:%s", err)
	}
	defer resp.Body.Close()
	return this.dealRestResponse(resp.Body)
}
