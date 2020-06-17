package ddxf_contract

import (
	"errors"
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology/common"
	"io"
)

type DataMetaInfo struct {
	DataMeta     map[string]interface{} `json:"dataMeta"`
	DataMetaHash string                 `json:"dataMetaHash"`
	ResourceType byte                   `json:"resourceType"`
	Fee          Fee                    `json:"fee"`
	Stock        uint32                 `json:"stock"`
	ExpiredDate  uint64                 `json:"expiredDate"`
	DataEndpoint string                 `json:"dataEndpoint"`
	DataHash     string                 `json:"dataHash"`
	DataId       string                 `json:"dataId"`
}

type TokenMetaInfo struct {
	TokenMeta     map[string]interface{} `json:"tokenMeta"`
	TokenMetaHash string                 `json:"tokenMetaHash"`
	DataMetaHash  string                 `json:"dataMetaHash"`
	TokenEndpoint string                 `json:"tokenEndpoint"`
}

type ProductInfoOnChain struct {
	ResourceDdo *ResourceDDO
	DtokenItem  *DTokenItem
}

type ResourceIdAndN struct {
	ResourceId []byte
	N          int
}

func (this *ProductInfoOnChain) Deserialize(source *common.ZeroCopySource) error {
	ddo := &ResourceDDO{}
	err := ddo.Deserialize(source)
	if err != nil {
		return err
	}
	item := &DTokenItem{}
	err = item.Deserialize(source)
	if err != nil {
		return err
	}
	this.ResourceDdo = ddo
	this.DtokenItem = item
	return nil
}

func (this *ProductInfoOnChain) FromBytes(data []byte) error {
	source := common.NewZeroCopySource(data)
	return this.Deserialize(source)
}

type CountAndAgent struct {
	Count  uint32
	Agents map[common.Address]uint32
}

func (this *CountAndAgent) FromBytes(data []byte) error {
	source := common.NewZeroCopySource(data)
	d, eof := source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}
	l, eof := source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}
	m := make(map[common.Address]uint32)
	for i := uint32(0); i < l; i++ {
		addr, eof := source.NextAddress()
		if eof {
			return io.ErrUnexpectedEOF
		}
		v, eof := source.NextUint32()
		if eof {
			return io.ErrUnexpectedEOF
		}
		m[addr] = v
	}
	this.Count = d
	this.Agents = m
	return nil
}

type TokenTemplate struct {
	DataID     string // can be empty
	TokenHashs []string
}

func (this *TokenTemplate) Deserialize(source *common.ZeroCopySource) error {
	data, irregular, eof := source.NextBool()
	if irregular || eof {
		return errors.New("")
	}
	if data {
		dataIds, _, irregular, eof := source.NextString()
		if irregular || eof {
			return fmt.Errorf("read dataids failed irregular:%v, eof:%v", irregular, eof)
		}
		this.DataID = dataIds
	}
	l, _, irregular, eof := source.NextVarUint()
	if irregular || eof {
		return fmt.Errorf("read tokenhash length failed irregular:%v, eof:%v", irregular, eof)
	}
	tokenHashs := make([]string, l)
	for i := 0; i < int(l); i++ {
		tokenHashs[i], _, irregular, eof = source.NextString()
		if irregular || eof {
			return fmt.Errorf("read tokenhash failed irregular:%v, eof:%v", irregular, eof)
		}
	}
	this.TokenHashs = tokenHashs
	return nil
}

func (this TokenTemplate) Serialize(sink *common.ZeroCopySink) {
	if len(this.DataID) == 0 {
		sink.WriteBool(false)
	} else {
		sink.WriteBool(true)
		sink.WriteString(this.DataID)
	}
	sink.WriteVarUint(uint64(len(this.TokenHashs)))
	for i := 0; i < len(this.TokenHashs); i++ {
		sink.WriteString(this.TokenHashs[i])
	}
}

func (this *TokenTemplate) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}

func (this *TokenTemplate) FromBytes(data []byte) error {
	source := common.NewZeroCopySource(data)
	return this.Deserialize(source)
}

type TokenResourceTyEndpoint struct {
	TokenTemplate *TokenTemplate
	ResourceType  byte
	Endpoint      string
}

func (this TokenResourceTyEndpoint) Serialize(sink *common.ZeroCopySink) {
	this.TokenTemplate.Serialize(sink)
	sink.WriteByte(this.ResourceType)
	sink.WriteString(this.Endpoint)
}
func (this *TokenResourceTyEndpoint) Deserialize(source *common.ZeroCopySource) error {
	err := this.TokenTemplate.Deserialize(source)
	if err != nil {
		return err
	}
	var eof bool
	this.ResourceType, eof = source.NextByte()
	if eof {
		return errors.New("read resource type failed")
	}
	return nil
}

// ResourceDDO is ddo for resource
type ResourceDDO struct {
	Manager                  common.Address             // data owner id
	TokenResourceTyEndpoints []*TokenResourceTyEndpoint // RT for tokens
	ItemMetaHash             common.Uint256             //
	DTC                      common.Address             // can be empty
	MP                       common.Address             // can be empty
	Split                    common.Address             // can be empty
	IsFreeze                 bool
}

func (this *ResourceDDO) Serialize(sink *common.ZeroCopySink) {
	sink.WriteAddress(this.Manager)
	sink.WriteVarUint(uint64(len(this.TokenResourceTyEndpoints)))
	for _, v := range this.TokenResourceTyEndpoints {
		v.Serialize(sink)
	}
	//TODO
	sink.WriteHash(this.ItemMetaHash)
	if this.DTC != common.ADDRESS_EMPTY {
		sink.WriteBool(true)
		sink.WriteAddress(this.DTC)
	} else {
		sink.WriteBool(false)
	}
	if this.MP != common.ADDRESS_EMPTY {
		sink.WriteBool(true)
		sink.WriteAddress(this.MP)
	} else {
		sink.WriteBool(false)
	}
	if this.Split != common.ADDRESS_EMPTY {
		sink.WriteBool(true)
		sink.WriteAddress(this.Split)
	} else {
		sink.WriteBool(false)
	}
	sink.WriteBool(this.IsFreeze)
}
func (this *ResourceDDO) Deserialize(source *common.ZeroCopySource) error {
	var eof bool
	this.Manager, eof = source.NextAddress()
	if eof {
		return io.ErrUnexpectedEOF
	}
	l, _, irregular, eof := source.NextVarUint()
	if irregular || eof {
		return errors.New("1. ResourceDDO Deserialize l error")
	}
	tokenResourceTypes := make([]*TokenResourceTyEndpoint, l)
	for i := 0; i < int(l); i++ {
		tt := &TokenResourceTyEndpoint{}
		err := tt.Deserialize(source)
		if err != nil {
			return err
		}
		tokenResourceTypes[i] = tt
	}
	this.TokenResourceTyEndpoints = tokenResourceTypes

	this.ItemMetaHash, eof = source.NextHash()
	if irregular || eof {
		return errors.New("2. ResourceDDO Deserialize l error")
	}
	data, irregular, eof := source.NextBool()
	if irregular || eof {
		return fmt.Errorf("read dtc failed irregular:%v, eof:%v", irregular, eof)
	}
	if data {
		this.DTC, eof = source.NextAddress()
		if eof {
			return io.ErrUnexpectedEOF
		}
	}
	data, irregular, eof = source.NextBool()
	if irregular || eof {
		return fmt.Errorf("read mp failed irregular:%v, eof:%v", irregular, eof)
	}
	if data {
		this.MP, eof = source.NextAddress()
		if eof {
			return io.ErrUnexpectedEOF
		}
	}
	data, irregular, eof = source.NextBool()
	if irregular || eof {
		return fmt.Errorf("read split failed irregular:%v, eof:%v", irregular, eof)
	}
	if data {
		this.Split, eof = source.NextAddress()
		if eof {
			return io.ErrUnexpectedEOF
		}
	}
	isFreeze, irregular, eof := source.NextBool()
	if irregular || eof {
		return fmt.Errorf("read isFreeze failed irregular:%v, eof:%v", irregular, eof)
	}
	if data {
		this.Split, eof = source.NextAddress()
		if eof {
			return io.ErrUnexpectedEOF
		}
	}
	this.IsFreeze = isFreeze
	return nil
}

func (this *ResourceDDO) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}

type Fee struct {
	ContractAddr common.Address
	ContractType split_policy_contract.TokenType
	Count        uint64
}

func (this *Fee) Serialize(sink *common.ZeroCopySink) {
	sink.WriteAddress(this.ContractAddr)
	sink.WriteByte(byte(this.ContractType))
	sink.WriteUint64(this.Count)
}
func (this *Fee) Deserialize(source *common.ZeroCopySource) error {
	var eof bool
	this.ContractAddr, eof = source.NextAddress()
	if eof {
		return io.ErrUnexpectedEOF
	}
	ty, eof := source.NextByte()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.ContractType = split_policy_contract.TokenType(ty)
	this.Count, eof = source.NextUint64()
	if eof {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type DTokenItem struct {
	Fee         Fee
	ExpiredDate uint64
	Stocks      uint32
	Sold        uint32
	Templates   []*TokenTemplate
}

func (this *DTokenItem) Serialize(sink *common.ZeroCopySink) {
	this.Fee.Serialize(sink)
	sink.WriteUint64(this.ExpiredDate)
	sink.WriteUint32(this.Stocks)
	sink.WriteUint32(this.Sold)
	sink.WriteVarUint(uint64(len(this.Templates)))
	for _, item := range this.Templates {
		item.Serialize(sink)
	}
}
func (this *DTokenItem) Deserialize(source *common.ZeroCopySource) error {
	err := this.Fee.Deserialize(source)
	if err != nil {
		return err
	}
	var eof bool
	this.ExpiredDate, eof = source.NextUint64()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.Stocks, eof = source.NextUint32()
	if eof {
		return fmt.Errorf("read stocks failed, eof: %v", eof)
	}
	this.Sold, eof = source.NextUint32()
	if eof {
		return fmt.Errorf("read sold failed, eof: %v", eof)
	}
	l, _, irre, eof := source.NextVarUint()
	if irre || eof {
		return fmt.Errorf("read tokentemplate length failed, irre: %v, eof: %v", irre, eof)
	}
	tts := make([]*TokenTemplate, l)
	for i := 0; i < int(l); i++ {
		tt := &TokenTemplate{}
		err = tt.Deserialize(source)
		if err != nil {
			return err
		}
		tts[i] = tt
	}
	return nil
}

func (this *DTokenItem) FromBytes(data []byte) error {
	source := common.NewZeroCopySource(data)
	return this.Deserialize(source)
}

func (this *DTokenItem) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}
