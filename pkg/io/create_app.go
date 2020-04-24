package io

import (
	"fmt"

	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// CreateAppInput for input
type CreateAppInput model.App

// Validate CreateAppInput
func (input *CreateAppInput) Validate() (err error) {
	switch {
	case input.ID == 0:
		err = fmt.Errorf("id empty")
	case input.Name == "":
		err = fmt.Errorf("name empty")

	}
	return
}

// CreateAppOutput for output
type CreateAppOutput struct {
	BaseResp
}
