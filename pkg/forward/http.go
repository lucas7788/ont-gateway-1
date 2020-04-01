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

// PostJSON for post a json http request
func PostJSON(uri string, data interface{}) (code int, respBody []byte, err error) {
	body, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", uri, strings.NewReader(string(body)))
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
