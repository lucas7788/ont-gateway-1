package client

import (
	"encoding/json"
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/registry/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
)

type RegistryImplClient struct {
	addr string
}

func Sdk() *RegistryImplClient {
	return &RegistryImplClient{
		addr: "http://127.0.0.1:20331",
	}
}

func (this *RegistryImplClient) AddEndpoint(input io.RegistryAddEndpointInput) (output io.RegistryAddEndpointOutput) {
	paramBs, err := json.Marshal(input)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	_, _, res, err := forward.PostJSONRequest(this.addr+server.AddEndpoint, paramBs, nil)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	err = json.Unmarshal(res, &output)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func (this *RegistryImplClient) RemoveEndpoint(input io.RegistryRemoveEndpointInput) (output io.RegistryRemoveEndpointOutput) {
	paramBs, err := json.Marshal(input)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	_, _, res, err := forward.PostJSONRequest(this.addr+server.RemoveEndpoint, paramBs, nil)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	err = json.Unmarshal(res, &output)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func (this *RegistryImplClient) QueryEndpoints(input io.RegistryQueryEndpointsInput) (output io.RegistryQueryEndpointsOutput) {
	_, _, bs, err := forward.Get(this.addr + server.QueryEndpoint)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	err = json.Unmarshal(bs, &output)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}
