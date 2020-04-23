package misc

import (
	osdk "github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
)

// OntSdk deals with ont sdk related stuff
type OntSdk struct {
	kit *osdk.OntologySdk
}

// NewOntSdk is ctor for OntSdk
func NewOntSdk() *OntSdk {

	sdk := &OntSdk{}
	sdk.init()

	return sdk
}

func (sdk *OntSdk) init() {
	kit := osdk.NewOntologySdk()
	{
		kit.NewRpcClient().SetAddress(sdk.GetOntNode())
	}
	sdk.kit = kit
}

// GetOntNode returns the ont node
func (sdk *OntSdk) GetOntNode() string {
	if config.Load().Prod {
		return "http://dappnode1.ont.io:20336"
	}

	return "http://polaris1.ont.io:20336"
}

// GetKit returns the ont sdk
func (sdk *OntSdk) GetKit() *osdk.OntologySdk {
	return sdk.kit
}

const (
	// ONGContractAddr for ONG contract address
	ONGContractAddr = "0200000000000000000000000000000000000000"
)

// GetAmountTransferred returns the amount transferred
func (sdk *OntSdk) GetAmountTransferred(txHash string) (amount uint64, err error) {

	event, err := sdk.kit.GetSmartContractEvent(txHash)
	if err != nil {
		return
	}

	for _, notify := range event.Notify {
		if notify.ContractAddress == ONGContractAddr {
			states, ok := notify.States.([]interface{})
			if !ok {
				continue
			}
			if len(states) >= 4 && states[0] == "transfer" {
				amount += states[3].(uint64)
			}
		}
	}

	return
}
