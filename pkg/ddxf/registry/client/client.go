package client

import (
	"encoding/json"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/registry/server"
	"net/http"
)

type RegistryImplClient struct {
	httpClient *HttpClient
}

func Sdk(addr string) *RegistryImplClient {
	return &RegistryImplClient{
		httpClient: NewHttpClient(addr),
	}
}

func (this *RegistryImplClient) AddEndpoint(input io.RegistryAddEndpointInput) (output io.RegistryAddEndpointOutput) {
	bs,err := this.httpClient.SendPostRequest(input,server.AddEndpoint)
	if err != nil {
		output.Code = http.StatusInternalServerError
		return
	}
	err = json.Unmarshal(bs, &output)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func (this *RegistryImplClient) RemoveEndpoint(input io.RegistryRemoveEndpointInput) (output io.RegistryRemoveEndpointOutput) {
	bs,err := this.httpClient.SendPostRequest(input,server.RemoveEndpoint)
	if err != nil {
		output.Code = http.StatusInternalServerError
		return
	}
	err = json.Unmarshal(bs, &output)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func (this *RegistryImplClient) QueryEndpoints(input io.RegistryQueryEndpointsInput) (output io.RegistryQueryEndpointsOutput) {
	bs, err := this.httpClient.SendGetRequest(server.QueryEndpoint)
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
