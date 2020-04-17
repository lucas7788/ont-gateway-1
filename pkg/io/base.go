package io

import "fmt"

// BaseResp for resp std
type BaseResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (r *BaseResp) Error() error {
	if r.Code == 0 || r.Code == 200 {
		return nil
	}

	return fmt.Errorf("code:%d msg:%s", r.Code, r.Msg)
}
