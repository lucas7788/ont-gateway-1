package misc

import (
	"github.com/ontio/ontology-go-sdk"
	common2 "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

type DDXFContractKit struct {
	bc              *BaseContract
	contractAddress common.Address
}

func NewDDXFContractKit(sdk *ontology_go_sdk.OntologySdk,
	contractAddress common.Address, gasLimit uint64,
	gasPrice uint64,
	payer *ontology_go_sdk.Account) *DDXFContractKit {
	bc := &BaseContract{
		sdk:      sdk,
		gasLimit: gasLimit,
		gasPrice: gasPrice,
		payer:    payer,
	}
	return &DDXFContractKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *DDXFContractKit) PreInvoke(method string, param []interface{}) (*common2.ResultItem, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, method, param)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *DDXFContractKit) Invoke(signer *ontology_go_sdk.Account, method string, param []interface{}) (common.Uint256, error) {
	txhash, err := this.bc.Invoke(this.contractAddress, signer, method, param)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return txhash, nil
}
func (this *DDXFContractKit) BuildTx(signer *ontology_go_sdk.Account, method string, param []interface{}) (*types.MutableTransaction, error) {
	tx, err := this.bc.BuildTx(this.contractAddress, signer, method, param)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
