package service

import "sync"

// Gateway service
type Gateway struct {
}

var (
	instanceGateway *Gateway
	lockGateway     sync.Mutex
)

// Instance is singleton for Gateway
func Instance() *Gateway {
	if instanceGateway != nil {
		return instanceGateway
	}

	lockGateway.Lock()
	defer lockGateway.Unlock()
	if instanceGateway != nil {
		return instanceGateway
	}

	instanceGateway = &Gateway{}

	return instanceGateway
}
