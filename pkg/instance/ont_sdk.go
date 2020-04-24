package instance

import (
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
)

var (
	ontSdk     *misc.OntSdk
	ontSdkLock sync.Mutex
)

// OntSdk is singleton for misc.OntSdk
func OntSdk() *misc.OntSdk {
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
