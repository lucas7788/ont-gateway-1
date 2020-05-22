package service

import (
	"net/http"

	"github.com/piprate/json-gold/ld"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// JSONLDAlign impl
func (gw *Gateway) JSONLDAlign(input io.JSONLDAlignInput) (output io.JSONLDAlignOutput) {

	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	expandA, err := proc.Expand(input.A, options)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	expandB, err := proc.Expand(input.B, options)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	jsonA := expandA[0].(map[string]interface{})
	jsonB := expandB[0].(map[string]interface{})
	jsonC := make(map[string]interface{})
	for k, v := range jsonA {
		if _, ok := jsonB[k]; ok {
			jsonC[k] = v
		}
	}

	ctxB := input.A["@context"]

	output.A2B, err = proc.Compact(jsonC, ctxB, options)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	return
}
