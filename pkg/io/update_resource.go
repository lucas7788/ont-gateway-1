package io

import "github.com/zhiqiangxu/ont-gateway/pkg/model"

// UpdateResourceInput ...
type UpdateResourceInput struct {
	RV    model.ResourceVersion `json:"rv"`
	Force bool                  `json:"force"`
}

// UpdateResourceOutput ...
type UpdateResourceOutput struct {
	BaseResp
	Exists bool `json:"exists"`
}
