package instance

import (
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
)

var (
	ontSdk     *misc.OntSdk
	ontSdkOnce sync.Once
)

// OntSdk is singleton for misc.OntSdk
func OntSdk() *misc.OntSdk {
	ontSdkOnce.Do(func() {
		ontSdk = misc.NewOntSdk()
	})

	return ontSdk
}
