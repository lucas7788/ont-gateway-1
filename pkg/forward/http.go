package forward

import (
	"encoding/json"
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

// PostJSONBytes for post marshaled json http request
func PostJSONBytes(uri string, data []byte) (code int, respBody []byte, err error) {
	req, err := http.NewRequest("POST", uri, strings.NewReader(string(data)))
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

	return
}

// PostJSON for post a json http request
func PostJSON(uri string, data interface{}) (code int, respBody []byte, err error) {
	body, err := json.Marshal(data)
	if err != nil {
		return
	}

	code, respBody, err = PostJSONBytes(uri, body)

	return
}
