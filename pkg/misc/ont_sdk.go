package misc

import (
	osdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	common2 "github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"time"
)

var ddxfContract *DDXFContractKit

const DDXF_CONTRACT_ADDRESS = "cf267b778d54174717d2fe81f2a931fcffc2cdd4"

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

func (sdk *OntSdk) DDXFContract(gasLimit uint64,
	gasPrice uint64,
	payer *osdk.Account) *DDXFContractKit {
	if ddxfContract == nil {
		contractAddress, _ := common2.AddressFromHexString(DDXF_CONTRACT_ADDRESS)
		ddxfContract = NewDDXFContractKit(sdk.kit, contractAddress)
	} else {
		ddxfContract.gasPrice = gasPrice
		ddxfContract.gasLimit = gasLimit
		ddxfContract.payer = payer
	}
	return ddxfContract
}

func (sdk *OntSdk) SendTx(txHex string) (string, error) {
	tx, err := utils.TransactionFromHexString(txHex)
	if err != nil {
		return "", err
	}
	mutTx, err := tx.IntoMutable()
	if err != nil {
		return "", err
	}
	txHash, err := instance.OntSdk().GetKit().SendTransaction(mutTx)
	return txHash.ToHexString(), err
}

func (sdk *OntSdk) GetSmartCodeEvent(txHash string) (*common.SmartContactEvent, error) {
	return sdk.kit.GetSmartContractEvent(txHash)
}

func (sdk *OntSdk) WaitForGenerateBlock() (bool, error) {
	timeout := time.Second * 60
	return sdk.kit.WaitForGenerateBlock(timeout)
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
