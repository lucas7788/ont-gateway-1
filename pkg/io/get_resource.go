package io

// GetResourceInput ...
type GetResourceInput struct {
	App   int    `json:"app"`
	ID    string `json:"id"`
	Block uint32 `json:"block"`
	Hash  string `json:"hash"`
}

// GetResourceOutput ...
type GetResourceOutput struct {
	BaseResp
	Exists   bool   `json:"exists"`
	Desc     string `json:"desc"`
	DescHash string `json:"desc_hash"`
}
