package service

import (
	"net/http"

	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"github.com/zhiqiangxu/util"
)

// UpdateResource ...
func (gw *Gateway) UpdateResource(input io.UpdateResourceInput) (output io.UpdateResourceOutput) {

	hash := ddxf.Sha256Bytes(util.Slice(input.RV.Desc))
	if string(hash[:]) != input.RV.Hash {
		output.Code = http.StatusBadRequest
		output.Msg = "bad hash"
		return
	}

	var err error
	if input.Force {

		output.Exists, err = model.ResourceVersionManager().ForceUpdateResource(input.RV)

	} else {

		output.Exists, err = model.ResourceVersionManager().UpdateResource(input.RV)

	}

	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}
