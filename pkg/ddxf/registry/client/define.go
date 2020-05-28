package client

import "encoding/json"

const REST_VERSION = "1.0.0"

type RestfulReq struct {
	Action  string
	Version string
	Type    int
	Data    string
}

type RestfulResp struct {
	Action  string          `json:"action"`
	Result  json.RawMessage `json:"result"`
	Error   int64           `json:"error"`
	Desc    string          `json:"desc"`
	Version string          `json:"version"`
}
