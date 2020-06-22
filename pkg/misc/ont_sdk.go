package misc

import (
	"time"

	"fmt"

	osdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	common2 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
)

var (
	ddxfContract        *DDXFContractKit
	defaultDdxfContract *DDXFContractKit

	dataIdContract        *DataIdContractKit
	defaultDataIdContract *DataIdContractKit
)

const DDXF_CONTRACT_ADDRESS = "90982cd1d33ec7b33bffe54b289f5acaf02815a8"
const DATA_ID_CONTRACT_ADDRESS = "e854316627dfc44bef9c0eb583e941804d0716d5"

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
	return GetOntNode()
}

// GetOntNode for direct usage
func GetOntNode() string {
	if config.Load().Prod {
		return "http://dappnode1.ont.io:20336"
	}

	return "http://polaris1.ont.io:20336"
}

// GetKit returns the ont sdk
func (sdk *OntSdk) GetKit() *osdk.OntologySdk {
	return sdk.kit
}

func (sdk *OntSdk) DefaultDataIdContract() *DataIdContractKit {
	if defaultDataIdContract == nil {
		contractAddress, _ := common2.AddressFromHexString(DATA_ID_CONTRACT_ADDRESS)
		defaultDataIdContract = NewDataIdContractKit(sdk.kit, contractAddress, 20000000, 500, nil)
	}
	return defaultDataIdContract
}

func (sdk *OntSdk) DataIdContract(gasLimit uint64,
	gasPrice uint64,
	payer *osdk.Account) *DataIdContractKit {
	if ddxfContract == nil {
		contractAddress, _ := common2.AddressFromHexString(DDXF_CONTRACT_ADDRESS)
		dataIdContract = NewDataIdContractKit(sdk.kit, contractAddress, gasLimit, gasPrice, payer)
	} else {
		dataIdContract.bc.gasPrice = gasPrice
		dataIdContract.bc.gasLimit = gasLimit
		dataIdContract.bc.payer = payer
	}
	return dataIdContract
}

func (sdk *OntSdk) DefaultDDXFContract() *DDXFContractKit {
	if defaultDdxfContract == nil {
		contractAddress, _ := common2.AddressFromHexString(DDXF_CONTRACT_ADDRESS)
		defaultDdxfContract = NewDDXFContractKit(sdk.kit, contractAddress, 20000000, 500, nil)
	}
	return defaultDdxfContract
}

func (sdk *OntSdk) DDXFContract(gasLimit uint64,
	gasPrice uint64,
	payer *osdk.Account) *DDXFContractKit {
	if ddxfContract == nil {
		contractAddress, _ := common2.AddressFromHexString(DDXF_CONTRACT_ADDRESS)
		ddxfContract = NewDDXFContractKit(sdk.kit, contractAddress, gasLimit, gasPrice, payer)
	} else {
		ddxfContract.bc.gasPrice = gasPrice
		ddxfContract.bc.gasLimit = gasLimit
		ddxfContract.bc.payer = payer
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
	txHash, err := sdk.kit.SendTransaction(mutTx)
	return txHash.ToHexString(), err
}
func (sdk *OntSdk) SendRawTx(mutTx *types.MutableTransaction) (string, error) {
	txHash, err := sdk.kit.SendTransaction(mutTx)
	return txHash.ToHexString(), err
}

func (sdk *OntSdk) GetSmartCodeEvent(txHash string) (*common.SmartContactEvent, error) {
	for i := 0; i < 10; i++ {
		event, err := sdk.kit.GetSmartContractEvent(txHash)
		if event != nil {
			return event, err
		}
		if err != nil {
			fmt.Println("GetSmartContractEvent err:", err)
			return nil, err
		}
		if event == nil {
			time.Sleep(3 * time.Second)
		}
	}
	return nil, fmt.Errorf("GetSmartCodeEvent timeout, txhash: %s", txHash)
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
