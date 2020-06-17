package instance

import (
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
)

var (
	ontSdk     *misc.OntSdk
	ontSdkOnce sync.Once
	ddxfSdk *ddxf_sdk.DdxfSdk
)

// OntSdk is singleton for misc.OntSdk
func OntSdk() *misc.OntSdk {
	ontSdkOnce.Do(func() {
		ontSdk = misc.NewOntSdk()
	})

	return ontSdk
}

// OntSdk is singleton for misc.OntSdk
func DDXFSdk() *ddxf_sdk.DdxfSdk {
	ontSdkOnce.Do(func() {
		if config.Load().Prod {
			ddxfSdk = ddxf_sdk.NewDdxfSdk("http://dappnode1.ont.io:20336")
		} else {
			ddxfSdk = ddxf_sdk.NewDdxfSdk("http://polaris1.ont.io:20336")
		}
	})
	return ddxfSdk
}