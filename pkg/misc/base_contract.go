package misc

import (
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	common2 "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

type BaseContract struct {
	sdk      *ontology_go_sdk.OntologySdk
	gasLimit uint64
	gasPrice uint64
	payer    *ontology_go_sdk.Account
}

func (this *BaseContract) PreInvoke(contractAddr common.Address, method string, param []interface{}) (*common2.ResultItem, error) {
	res, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(contractAddr, method, param)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (this *BaseContract) Invoke(contractAddr common.Address, signer *ontology_go_sdk.Account, method string, param []interface{}) (common.Uint256, error) {
	txhash, err := this.sdk.WasmVM.InvokeWasmVMSmartContract(this.gasPrice, this.gasLimit, this.payer, signer, contractAddr, method, param)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return txhash, nil
}
func (this *BaseContract) BuildTx(contractAddr common.Address, signer *ontology_go_sdk.Account, method string, param []interface{}) (*types.MutableTransaction, error) {
	tx, err := this.sdk.WasmVM.NewInvokeWasmVmTransaction(this.gasPrice, this.gasLimit, contractAddr, method, param)
	if err != nil {
		return nil, err
	}
	if this.payer != nil {
		this.sdk.SetPayer(tx, this.payer.Address)
		err = this.sdk.SignToTransaction(tx, this.payer)
		if err != nil {
			return nil, fmt.Errorf("payer sign tx error: %s", err)
		}
	}
	if this.payer != signer {
		err = this.sdk.SignToTransaction(tx, signer)
		if err != nil {
			return nil, err
		}
	}
	return tx, nil
}
