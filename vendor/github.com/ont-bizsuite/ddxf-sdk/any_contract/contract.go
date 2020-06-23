package any_contract

import (
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ontio/ontology-go-sdk"
	common2 "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

type ContractKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewContractKit(contractAddress common.Address, bc *base_contract.BaseContract) *ContractKit {
	return &ContractKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *ContractKit) SetContractAddress(addr common.Address) {
	this.contractAddress = addr
}

func (this *ContractKit) Invoke(method string, signer *ontology_go_sdk.Account, args []interface{}) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, signer, method, args)
}

func (this *ContractKit) BuildTx(method string, signer *ontology_go_sdk.Account, args []interface{}) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, method, args)
}

func (this *ContractKit) PreInvoke(method string, args []interface{}) (*common2.ResultItem, error) {
	return this.bc.PreInvoke(this.contractAddress, method, args)
}
