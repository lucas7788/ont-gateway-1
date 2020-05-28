package misc

import (
	"github.com/ontio/ontology-go-sdk"
	common2 "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
)

type DDXFContractKit struct {
	sdk             *ontology_go_sdk.OntologySdk
	gasLimit        uint64
	gasPrice        uint64
	contractAddress common.Address
	payer           *ontology_go_sdk.Account
}

func NewDDXFContractKit(sdk *ontology_go_sdk.OntologySdk,
	contractAddress common.Address) *DDXFContractKit {
	return &DDXFContractKit{
		sdk:             sdk,
		contractAddress: contractAddress,
	}
}

func (this *DDXFContractKit) PreInvoke(method string, param []interface{}) (*common2.ResultItem, error) {
	res, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(this.contractAddress, method, param)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (this *DDXFContractKit) Invoke(signer *ontology_go_sdk.Account, method string, param []interface{}) (common.Uint256, error) {
	txhash, err := this.sdk.WasmVM.InvokeWasmVMSmartContract(this.gasPrice, this.gasLimit, this.payer, signer, this.contractAddress, method, param)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return txhash, nil
}
