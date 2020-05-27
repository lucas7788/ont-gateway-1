package registry

import "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"

// Registry ...
type Registry interface {
	AddEndpoint(io.RegistryAddEndpointInput) io.RegistryAddEndpointOutput
	RemoveEndpoint(io.RegistryRemoveEndpointOutput) io.RegistryRemoveEndpointOutput
	QueryEndpoints(io.RegistryQueryEndpointsInput) io.RegistryQueryEndpointsOutput
	Sdk() Registry
}
