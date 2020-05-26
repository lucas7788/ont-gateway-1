package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

// RegistryQueryMPEndpointsInput ...
type RegistryQueryMPEndpointsInput struct {
}

// MPEndpoint ...
type MPEndpoint struct {
	MP       string
	Endpoint string
}

// RegistryQueryMPEndpointsOutput ...
type RegistryQueryMPEndpointsOutput struct {
	io2.BaseResp
	Endpoints []MPEndpoint
}
