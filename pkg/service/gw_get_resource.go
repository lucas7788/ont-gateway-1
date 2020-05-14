package service

import (
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// GetResource ...
func (gw *Gateway) GetResource(input io.GetResourceInput) (output io.GetResourceOutput) {

	var err error
	var rv *model.ResourceVersion
	if input.Block != 0 {
		rv, err = model.ResourceVersionManager().GetByBlock(input.App, input.ID, input.Block)
	} else {
		rv, err = model.ResourceVersionManager().GetByHash(input.App, input.ID, input.Hash)
	}

	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	output.Desc = rv.Desc
	output.DescHash = rv.Hash
	return
}
