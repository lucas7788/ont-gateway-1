package client

import "encoding/json"

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
