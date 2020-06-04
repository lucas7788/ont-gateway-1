package misc

import (
	"github.com/ontio/ontology-go-sdk"
	common3 "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

type DataIdContractKit struct {
	bc              *BaseContract
	contractAddress common.Address
}

func NewDataIdContractKit(sdk *ontology_go_sdk.OntologySdk,
	contractAddress common.Address, gasLimit uint64,
	gasPrice uint64,
	payer *ontology_go_sdk.Account) *DataIdContractKit {
	bc := &BaseContract{
		sdk:      sdk,
		gasLimit: gasLimit,
		gasPrice: gasPrice,
		payer:    payer,
	}
	return &DataIdContractKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *DataIdContractKit) PreInvoke(method string, param []interface{}) (*common3.ResultItem, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, method, param)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *DataIdContractKit) Invoke(signer *ontology_go_sdk.Account, method string, param []interface{}) (common.Uint256, error) {
	txhash, err := this.bc.Invoke(this.contractAddress, signer, method, param)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return txhash, nil
}
func (this *DataIdContractKit) BuildTx(signer *ontology_go_sdk.Account, method string, param []interface{}) (*types.MutableTransaction, error) {
	tx, err := this.bc.BuildTx(this.contractAddress, signer, method, param)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
