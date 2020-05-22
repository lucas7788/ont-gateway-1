package io

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// ListAppOutput for output
type ListAppOutput struct {
	BaseResp
	Apps []*model.App
}
