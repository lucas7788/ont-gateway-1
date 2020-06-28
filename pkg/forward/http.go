package forward

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/metrics"
)

const (
	// defaultTimeout is default timeout for api request
	defaultTimeout = 30 * time.Second
)

var (
	client = &http.Client{Timeout: defaultTimeout, Transport: &http.Transport{IdleConnTimeout: time.Second * 2, MaxIdleConnsPerHost: 200}}
)

func httpRequest(req *http.Request) (code int, contentType string, respBody []byte, err error) {
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	code = resp.StatusCode

	contentType = resp.Header.Get("Content-Type")
	return
}

// JSONRequest for send json request
func JSONRequest(method, uri string, data []byte, headers map[string]string) (code int, contentType string, respBody []byte, err error) {
	req, err := http.NewRequest(method, uri, strings.NewReader(string(data)))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	code, contentType, respBody, err = httpRequest(req)
	return
}

// PostJSONRequestWithRetry with retry-ability
func PostJSONRequestWithRetry(uri string, data []byte, headers map[string]string, n int) (code int, contentType string, respBody []byte, err error) {
	for i := 0; i < n; i++ {
		code, contentType, respBody, err = PostJSONRequest(uri, data, headers)
		if err == nil {
			return
		}
		if errNet, ok := err.(net.Error); ok && errNet.Timeout() {
			continue
		}
		return
	}
	return
}

// PostJSONRequest for post json request
func PostJSONRequest(uri string, data []byte, headers map[string]string) (int, string, []byte, error) {
	return JSONRequest("POST", uri, data, headers)
}

// Get request
func Get(url string) (code int, contentType string, respBody []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	code = resp.StatusCode

	contentType = resp.Header.Get("Content-Type")
	return
}

var (
	hopHeaders = map[string]bool{
		"Connection":          true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"Te":                  true, // canonicalized version of "TE"
		"Trailers":            true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}
)

// GenericForward forwards an arbitory request
func GenericForward(req *http.Request, prefix, remoteAddr, schema string) (resp *http.Response, err error) {
	path := req.URL.Path
	if req.URL.RawQuery != "" {
		path += "?" + req.URL.RawQuery
	}
	httpReq, err := http.NewRequest(req.Method, strings.TrimPrefix(path, prefix), req.Body)
	if err != nil {
		return
	}

	httpReq.URL.Host = remoteAddr
	httpReq.URL.Scheme = schema

	for header, values := range req.Header {
		if !hopHeaders[header] {
			for _, value := range values {
				httpReq.Header.Set(header, value)
			}
		}
	}

	begin := time.Now()
	resp, err = client.Do(httpReq)

	errStr := fmt.Sprintf("%v", err)
	lvs := []string{"method", "forward", "error", errStr}
	metrics.RequestLatencyMetric.With(lvs...).Observe(time.Since(begin).Seconds())

	return
}
