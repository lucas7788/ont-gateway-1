package market_place_contract

import (
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

type MpKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewDDXFContractKit(contractAddress common.Address, bc *base_contract.BaseContract) *MpKit {
	return &MpKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *MpKit) SetContractAddress(addr common.Address) {
	this.contractAddress = addr
}

func (this *MpKit) Init(admin *ontology_go_sdk.Account, dtoken, splitPolicy common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "init",
		[]interface{}{dtoken, splitPolicy})
}

func (this *MpKit) setDTokenContractAddress(admin *ontology_go_sdk.Account, dtoken common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "setDTokenContract", []interface{}{dtoken})
}

func (this *MpKit) getDTokenContractAddress() (common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getDTokenContract", []interface{}{})
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	return common.AddressParseFromBytes(data)
}

func (this *MpKit) setSplitPolicyContractAddress(admin *ontology_go_sdk.Account, dtoken common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "setSplitPolicyContract", []interface{}{dtoken})
}

func (this *MpKit) getSplitPolicyContractAddress() (common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getSplitPolicyContract", []interface{}{})
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	return common.AddressParseFromBytes(data)
}

//publish product on block chain,
func (this *MpKit) Publish(seller *ontology_go_sdk.Account, resourceId []byte, ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, seller, "dtokenSellerPublish",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
}

func (this *MpKit) Update(seller *ontology_go_sdk.Account, resourceId []byte, ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, seller, "update",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
}

func (this *MpKit) BuildUpdateTx(resourceId []byte, ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "update",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
}

func (this *MpKit) Delete(seller *ontology_go_sdk.Account, resourceId []byte) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, seller, "delete",
		[]interface{}{resourceId})
}

func (this *MpKit) BuildDeleteTx(resourceId []byte) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "delete",
		[]interface{}{resourceId})
}

func (this *MpKit) FreezeAndPublish(seller *ontology_go_sdk.Account, resourceIdOld, resourceIdNew []byte,
	ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, seller, "freezeAndPublish",
		[]interface{}{resourceIdOld, resourceIdNew, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
}

func (this *MpKit) BuildFreezeAndPublishTx(resourceIdOld, resourceIdNew []byte,
	ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "freezeAndPublish", []interface{}{
		resourceIdOld, resourceIdNew, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes(),
	})
}

func (this *MpKit) BuildPublishTx(resourceId []byte,
	ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (*types.MutableTransaction, error) {
	tx, err := this.bc.BuildTx(this.contractAddress, "dtokenSellerPublish",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
	return tx, err
}

func (this *MpKit) getPublishProductInfo(resourceId []byte) (*ProductInfoOnChain, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getSellerItemInfo",
		[]interface{}{resourceId})
	if err != nil {
		return nil, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return nil, err
	}
	p := &ProductInfoOnChain{}
	err = p.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (this *MpKit) Freeze(manager *ontology_go_sdk.Account,
	resourceId []byte) (common.Uint256, error) {
	tx, err := this.bc.BuildTx(this.contractAddress, "freeze",
		[]interface{}{resourceId})
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.bc.GetOntologySdk().SignToTransaction(tx, manager)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.bc.GetOntologySdk().SendTransaction(tx)
}

func (this *MpKit) BuildFreezeTx(resourceId []byte) (*types.MutableTransaction, error) {
	tx, err := this.bc.BuildTx(this.contractAddress, "freeze",
		[]interface{}{resourceId})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (this *MpKit) BuyDtoken(buyer, payer *ontology_go_sdk.Account, resourceId []byte,
	n int) (common.Uint256, error) {
	if payer == nil {
		payer = buyer
	}
	tx, err := this.bc.BuildTx(this.contractAddress, "buyDtoken",
		[]interface{}{resourceId, n, buyer.Address, payer.Address})
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.bc.GetOntologySdk().SignToTransaction(tx, buyer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if buyer.Address != payer.Address {
		err = this.bc.GetOntologySdk().SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	return this.bc.GetOntologySdk().SendTransaction(tx)
}

func (this *MpKit) BuyDtokenReward(buyer, payer *ontology_go_sdk.Account, resourceId []byte,
	n int, unitPrice int) (common.Uint256, error) {
	if payer == nil {
		payer = buyer
	}
	tx, err := this.bc.BuildTx(this.contractAddress, "buyDtoken",
		[]interface{}{resourceId, n, buyer.Address, payer.Address, unitPrice})
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if buyer.Address != payer.Address {
		err = this.bc.GetOntologySdk().SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	return this.bc.GetOntologySdk().SendTransaction(tx)
}

func (this *MpKit) BuyAndUseToken(buyer, payer *ontology_go_sdk.Account, resourceId []byte,
	n int, tokenTemplate TokenTemplate) (common.Uint256, error) {
	if payer == nil {
		payer = buyer
	}
	tx, err := this.bc.BuildTx(this.contractAddress, "buyAndUseToken",
		[]interface{}{resourceId, n, buyer.Address, payer.Address, tokenTemplate.ToBytes()})
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if buyer.Address != payer.Address {
		err = this.bc.GetOntologySdk().SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	return this.bc.GetOntologySdk().SendTransaction(tx)
}

func (this *MpKit) BuildBuyAndUseTokenTx(buyer, payer common.Address, resourceId []byte,
	n int, tokenTemplate TokenTemplate) (*types.MutableTransaction, error) {
	tx, err := this.bc.BuildTx(this.contractAddress, "buyAndUseToken",
		[]interface{}{resourceId, n, buyer, payer, tokenTemplate.ToBytes()})
	if err != nil {
		return nil, err
	}
	return tx, nil
}
