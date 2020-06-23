package dtoken_contract

import (
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

type DTokenKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewDTokenKit(contractAddress common.Address, bc *base_contract.BaseContract) *DTokenKit {
	return &DTokenKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *DTokenKit) SetContractAddr(addr common.Address) {
	this.contractAddress = addr
}

func (this *DTokenKit) SetMpContractAddr(admin *ontology_go_sdk.Account, addr common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "setDdxfContract",
		[]interface{}{addr})
}

func (this *DTokenKit) GetDDXFContractAddr() (common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getDdxfContract", []interface{}{})
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	bs, err := res.ToByteArray()
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	source := common.NewZeroCopySource(bs)
	addr, eof := source.NextAddress()
	if eof {
		return common.ADDRESS_EMPTY, fmt.Errorf("read address failed, eof: %v", eof)
	}
	return addr, nil
}

func (this *DTokenKit) UseToken(buyer *ontology_go_sdk.Account,
	tokenTemplate market_place_contract.TokenTemplate, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, buyer, "useToken",
		[]interface{}{buyer.Address, tokenTemplate.ToBytes(), n})
}

func (this *DTokenKit) UseTokenByAgents(tokenOwner common.Address,
	agent *ontology_go_sdk.Account, tokenTemplate market_place_contract.TokenTemplate, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, agent, "useTokenByAgent",
		[]interface{}{tokenOwner, agent.Address, tokenTemplate.ToBytes(), n})
}

func (this *DTokenKit) BuyDtokens(buyer *ontology_go_sdk.Account,
	param []market_place_contract.ResourceIdAndN) (common.Uint256, error) {
	rids := make([]interface{}, len(param))
	for i := 0; i < len(param); i++ {
		rids[i] = param[i].ResourceId
	}
	ns := make([]interface{}, len(param))
	for i := 0; i < len(param); i++ {
		ns[i] = param[i].N
	}
	return this.bc.Invoke(this.contractAddress, buyer, "buyDtokens",
		[]interface{}{rids, ns, buyer.Address})
}

func (this *DTokenKit) SetAgents(account *ontology_go_sdk.Account,
	agents []common.Address, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "setAgents",
		[]interface{}{account.Address, agents, n})
}

func (this *DTokenKit) SetTokenAgents(account *ontology_go_sdk.Account,
	agents []common.Address, tokenTemplate market_place_contract.TokenTemplate, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "setTokenAgents",
		[]interface{}{account.Address, agents, tokenTemplate.ToBytes(), n})
}

func (this *DTokenKit) AddAgents(account *ontology_go_sdk.Account,
	agents []common.Address, n int, templates []market_place_contract.TokenTemplate) (common.Uint256, error) {
	sink := common.NewZeroCopySink(nil)
	sink.WriteVarUint(uint64(len(templates)))
	for _,template := range templates {
		template.Serialize(sink)
	}
	return this.bc.Invoke(this.contractAddress, account, "addAgents",
		[]interface{}{account.Address, parseAddressArr(agents), n, sink.Bytes()})
}

func parseAddressArr(addrs []common.Address) []interface{} {
	agentArr := make([]interface{}, len(addrs))
	for i := 0; i < len(addrs); i++ {
		agentArr[i] = addrs[i]
	}
	return agentArr
}

func (this *DTokenKit) AddTokenAgents(account *ontology_go_sdk.Account,
	agents []common.Address, tokenTemplate market_place_contract.TokenTemplate, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "addTokenAgents",
		[]interface{}{account.Address, tokenTemplate.ToBytes(), parseAddressArr(agents), n})
}

func (this *DTokenKit) RemoveAgents(account *ontology_go_sdk.Account,
	agents []common.Address, templates []market_place_contract.TokenTemplate) (common.Uint256, error) {
	sink := common.NewZeroCopySink(nil)
	sink.WriteVarUint(uint64(len(templates)))
	for _,template := range templates {
		template.Serialize(sink)
	}
	return this.bc.Invoke(this.contractAddress, account, "removeAgents",
		[]interface{}{account.Address, parseAddressArr(agents), sink.Bytes()})
}

func (this *DTokenKit) RemoveTokenAgents(tokenTemplate market_place_contract.TokenTemplate, account *ontology_go_sdk.Account,
	agents []common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "removeTokenAgents",
		[]interface{}{account.Address,tokenTemplate.ToBytes(),parseAddressArr(agents)})
}
