package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// RegistryQueryEndpointsInput ...
type RegistryQueryEndpointsInput struct {
}

// MPEndpoint ...
type MPEndpoint struct {
	MP       string
	Endpoint string
}

// RegistryQueryEndpointsOutput ...
type RegistryQueryEndpointsOutput struct {
	io2.BaseResp
	Endpoints []MPEndpoint
}
