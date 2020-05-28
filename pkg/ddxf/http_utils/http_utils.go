package http_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClient struct {
	addr       string
	httpClient *http.Client
}

func NewHttpClient(addr string) *HttpClient {
	return &HttpClient{
		addr: addr,
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

func (this *HttpClient) getRequestUrl(reqPath string, values ...*url.Values) (string, error) {
	if !strings.HasPrefix(this.addr, "http") {
		this.addr = "http://" + this.addr
	}
	reqUrl, err := new(url.URL).Parse(this.addr)
	if err != nil {
		return "", fmt.Errorf("Parse address:%s error:%s", this.addr, err)
	}
	reqUrl.Path = reqPath
	if len(values) > 0 && values[0] != nil {
		reqUrl.RawQuery = values[0].Encode()
	}
	return reqUrl.String(), nil
}

func (this *HttpClient) SendPostRequest(reqParam interface{}, reqPath string, values ...*url.Values) ([]byte, error) {
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
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http body error:%s", err)
	}
	return data,nil
}

func (this *HttpClient) PostRequest(reqParam interface{}, reqUrl string) ([]byte, error) {
	reqData, err := json.Marshal(reqParam)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal error:%s", err)
	}
	resp, err := this.httpClient.Post(reqUrl, "application/json", bytes.NewReader(reqData))
	if err != nil {
		return nil, fmt.Errorf("send http post request error:%s", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http body error:%s", err)
	}
	return data,nil
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
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http body error:%s", err)
	}
	return data,nil
}
