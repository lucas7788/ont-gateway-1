package registry

import "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"

// Registry ...
type Registry interface {
	QueryMPEndpoints(io.RegistryQueryMPEndpointsInput) io.RegistryQueryMPEndpointsOutput
}
