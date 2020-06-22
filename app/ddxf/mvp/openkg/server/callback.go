package server

import (
	"encoding/json"

	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
)

func callback(output interface{}) {
	outBytes, _ := json.Marshal(output)

	if config.OpenKGCallbackURI == "" {
		return
	}

	code, _, respBytes, err := forward.PostJSONRequest(config.OpenKGCallbackURI, outBytes, nil)
	if code != 200 {
		instance.Logger().Error("openkg callback", zap.Int("code", code), zap.String("resp", string(respBytes)))
	}
	if err != nil {
		instance.Logger().Error("openkg callback", zap.Error(err))
	}
}
