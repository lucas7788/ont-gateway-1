package instance

import (
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
)

var (
	ontSdk     *misc.OntSdk
	ontSdkLock sync.Mutex
)

// OntSdkInstance is singleton for OntSdk
func OntSdkInstance() *misc.OntSdk {
	if ontSdk != nil {
		return ontSdk
	}

	ontSdkLock.Lock()
	defer ontSdkLock.Unlock()

	if ontSdk != nil {
		return ontSdk
	}

	ontSdk = misc.NewOntSdk()

	return ontSdk
}
