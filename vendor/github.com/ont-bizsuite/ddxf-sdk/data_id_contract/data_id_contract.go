package data_id_contract

import (
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

type DataIdKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewDataIdContractKit(
	contractAddress common.Address,
	bc *base_contract.BaseContract) *DataIdKit {
	return &DataIdKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *DataIdKit) SetContractAddress(dataId common.Address) {
	this.contractAddress = dataId
}

func (this *DataIdKit) RegisterDataIdInfo(info DataIdInfo,
	seller *ontology_go_sdk.Account) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, seller, "registerDataId",
		[]interface{}{info.ToBytes()})
}

func (this *DataIdKit) RegisterDataIdInfoArray(info []DataIdInfo,
	seller *ontology_go_sdk.Account) (common.Uint256, error) {
	param := make([]interface{}, len(info))
	for i := 0; i < len(info); i++ {
		param[i] = info[i].ToBytes()
	}
	return this.bc.Invoke(this.contractAddress, seller, "registerDataIdArray",
		[]interface{}{param})
}

func (this *DataIdKit) BuildRegisterDataIdInfoArrayTx(info []DataIdInfo) (*types.MutableTransaction, error) {
	param := make([]interface{}, len(info))
	for i := 0; i < len(info); i++ {
		param[i] = info[i].ToBytes()
	}
	return this.bc.BuildTx(this.contractAddress, "registerDataIdArray", []interface{}{param})
}

func (this *DataIdKit) GetDataIdInfo(dataId string) (*DataIdInfo, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getDataIdInfo",
		[]interface{}{dataId})
	if err != nil {
		return nil, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return nil, err
	}
	info := &DataIdInfo{}
	err = info.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return info, nil
}
