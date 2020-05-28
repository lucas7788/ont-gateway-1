package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// RegistryQueryEndpointsInput ...
type RegistryQueryEndpointsInput struct {
}

// MPEndpoint ...
type MPEndpoint struct {
	MP       string `json:"mp"`
	Endpoint string `json:"endpoint"`
}

// RegistryQueryEndpointsOutput ...
type RegistryQueryEndpointsOutput struct {
	io2.BaseResp
	Endpoints []MPEndpoint `json:"endpoints"`
}
