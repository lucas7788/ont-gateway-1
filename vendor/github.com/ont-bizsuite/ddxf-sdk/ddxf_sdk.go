package ddxf_sdk

import (
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk/any_contract"
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ont-bizsuite/ddxf-sdk/dtoken_contract"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ontio/ontology-go-sdk"
	common2 "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"time"
)

const (
	mpContractAddressTest     = "9d0203fc1c1a5019c53fdf62ae3232f5a72f5d80"
	dtokenContractAddressTest = "466b94488bf2ad1b1eec0ae7e49e40708e71a35d"

	mpContractAddress     = "e01d500ed0c1719b7750367ae59b4b2d308d1ceb"
	dtokenContractAddress = "466b94488bf2ad1b1eec0ae7e49e40708e71a35d"
)

const (
	MainNet  = "http://dappnode1.ont.io:20336"
	TestNet  = "http://polaris1.ont.io:20336"
	LocalNet = "http://127.0.0.1:20336"
)

const (
	defaultGasPrice = 500
	defaultGasLimit = 31600000
)

type DdxfSdk struct {
	sdk          *ontology_go_sdk.OntologySdk
	bc           *base_contract.BaseContract
	rpc          string
	defDDXFKit   *market_place_contract.MpKit
	defDTokenKit *dtoken_contract.DTokenKit
	gasPrice     uint64
	gasLimit     uint64
	payer        *ontology_go_sdk.Account
}

func NewDdxfSdk(addr string) *DdxfSdk {
	sdk := ontology_go_sdk.NewOntologySdk()
	sdk.NewRpcClient().SetAddress(addr)
	return &DdxfSdk{
		sdk:      sdk,
		rpc:      addr,
		gasPrice: defaultGasPrice,
		gasLimit: defaultGasLimit,
		bc:       base_contract.NewBaseContract(sdk, defaultGasLimit, defaultGasPrice, nil),
	}
}

func (sdk *DdxfSdk) DefContract(contractAddr common.Address) *any_contract.ContractKit {
	return any_contract.NewContractKit(contractAddr, sdk.bc)
}

func (sdk *DdxfSdk) SetPayer(payer *ontology_go_sdk.Account) {
	sdk.payer = payer
	sdk.bc.SetPayer(payer)
}

func (sdk *DdxfSdk) SetGasLimit(gasLimit uint64) {
	sdk.gasLimit = gasLimit
	sdk.bc.SetGasLimit(gasLimit)
}

func (sdk *DdxfSdk) SetGasPrice(gasPrice uint64) {
	sdk.gasPrice = gasPrice
	sdk.bc.SetGasPrice(gasPrice)
}

func (sdk *DdxfSdk) GetOntologySdk() *ontology_go_sdk.OntologySdk {
	return sdk.sdk
}

func (sdk *DdxfSdk) SignTx(tx *types.MutableTransaction, signer *ontology_go_sdk.Account) error {
	if sdk.payer != nil {
		sdk.sdk.SetPayer(tx, sdk.payer.Address)
		err := sdk.sdk.SignToTransaction(tx, sdk.payer)
		if err != nil {
			return fmt.Errorf("payer sign tx error: %s", err)
		}
	}
	if sdk.payer != signer {
		err := sdk.sdk.SignToTransaction(tx, signer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sdk *DdxfSdk) SendTx(tx *types.MutableTransaction) (common.Uint256, error) {
	return sdk.sdk.SendTransaction(tx)
}

func (sdk *DdxfSdk) DefDTokenKit() *dtoken_contract.DTokenKit {
	if sdk.defDTokenKit == nil {
		contractAddress, _ := common.AddressFromHexString(dtokenContractAddress)
		sdk.defDTokenKit = dtoken_contract.NewDTokenKit(contractAddress,
			sdk.bc)
	}
	return sdk.defDTokenKit
}

func (sdk *DdxfSdk) DefMpKit() *market_place_contract.MpKit {
	if sdk.defDDXFKit == nil {
		contractAddress, _ := common.AddressFromHexString(mpContractAddress)
		sdk.defDDXFKit = market_place_contract.NewDDXFContractKit(contractAddress,
			sdk.bc)
	}
	return sdk.defDDXFKit
}

func (sdk *DdxfSdk) SetMpContractAddress(ddxf common.Address) {
	sdk.DefMpKit().SetContractAddress(ddxf)
}

func (sdk *DdxfSdk) GetSmartCodeEvent(txHash string) (*common2.SmartContactEvent, error) {
	for i := 0; i < 10; i++ {
		event, err := sdk.sdk.GetSmartContractEvent(txHash)
		if event != nil {
			return event, err
		}
		if err != nil {
			return nil, err
		}
		if event == nil {
			time.Sleep(3 * time.Second)
		}
	}
	return nil, fmt.Errorf("GetSmartCodeEvent timeout, txhash: %s", txHash)
}

func (this *DdxfSdk) DeployContract(signer *ontology_go_sdk.Account, code,
	name, version, author, email, desc string) (common.Uint256, error) {
	return this.sdk.WasmVM.DeployWasmVMSmartContract(this.gasPrice,
		this.gasLimit, signer, code, name,
		version, author, email, desc)
}
