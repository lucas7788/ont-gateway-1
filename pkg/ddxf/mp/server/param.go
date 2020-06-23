package server

import "github.com/zhiqiangxu/ont-gateway/pkg/io"

type DeleteInput struct {
	SignedTx string
}

type DeleteOutput struct {
	io.BaseResp
}

type UpdateInput struct {
	SignedTx string
}

type UpdateOutput struct {
	io.BaseResp
}
