package forward

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	// defaultTimeout is default timeout for api request
	defaultTimeout = 5 * time.Second
)

var (
	client = &http.Client{Timeout: defaultTimeout, Transport: &http.Transport{IdleConnTimeout: time.Second * 2, MaxIdleConnsPerHost: 200}}
)

// JSONRequest for send json request
func JSONRequest(method, uri string, data []byte) (code int, contentType string, respBody []byte, err error) {
	req, err := http.NewRequest(method, uri, strings.NewReader(string(data)))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
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

// PostJSONRequest for post json request
func PostJSONRequest(uri string, data []byte) (int, string, []byte, error) {
	return JSONRequest("POST", uri, data)
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
