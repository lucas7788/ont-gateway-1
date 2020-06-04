package param

import (
	"fmt"
	"github.com/kataras/go-errors"
	"github.com/ontio/ontology/common"
	"io"
)

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
	DataIDs    string // can be empty
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
		this.DataIDs = dataIds
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

func (this *TokenTemplate) Serialize(sink *common.ZeroCopySink) {
	if len(this.DataIDs) == 0 {
		sink.WriteBool(false)
	} else {
		sink.WriteBool(true)
		sink.WriteString(this.DataIDs)
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

type TokenTemplateEndPoint struct {
	TokenTemplate TokenTemplate
	Endpoint      string
}

func (this TokenTemplateEndPoint) Serialize(sink *common.ZeroCopySink) {
	this.TokenTemplate.Serialize(sink)
	sink.WriteString(this.Endpoint)
}

func (this *TokenTemplateEndPoint) Deserialize(source *common.ZeroCopySource) error {
	this.TokenTemplate.Deserialize(source)
	var irregular, eof bool
	this.Endpoint, _, irregular, eof = source.NextString()
	if irregular || eof {
		return fmt.Errorf("[TokenTemplateEndPoint] read endpoint failed,irregular:%v, eof:%v", irregular, eof)
	}
	return nil
}

type TokenResourceType struct {
	TokenTemplate TokenTemplate
	ResourceType  byte
}

func (this TokenResourceType) Serialize(sink *common.ZeroCopySink) {
	this.TokenTemplate.Serialize(sink)
	sink.WriteByte(this.ResourceType)
}
func (this *TokenResourceType) Deserialize(source *common.ZeroCopySource) error {
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
	Manager            common.Address           // data owner id
	TokenResourceTypes []*TokenResourceType     // RT for tokens
	TokenEndpoints     []*TokenTemplateEndPoint // endpoint for tokens
	ItemMetaHash       common.Uint256           //
	DTC                common.Address           // can be empty
	MP                 common.Address           // can be empty
	Split              common.Address           // can be empty
}

func (this *ResourceDDO) Serialize(sink *common.ZeroCopySink) {
	sink.WriteAddress(this.Manager)
	sink.WriteVarUint(uint64(len(this.TokenResourceTypes)))
	for _, v := range this.TokenResourceTypes {
		v.Serialize(sink)
	}
	sink.WriteVarUint(uint64(len(this.TokenEndpoints)))
	for _, v := range this.TokenEndpoints {
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
	tokenResourceTypes := make([]*TokenResourceType, l)
	for i := 0; i < int(l); i++ {
		tt := &TokenResourceType{}
		err := tt.Deserialize(source)
		if err != nil {
			return err
		}
		tokenResourceTypes[i] = tt
	}
	this.TokenResourceTypes = tokenResourceTypes
	l, _, irregular, eof = source.NextVarUint()
	if irregular || eof {
		return errors.New("1. ResourceDDO Deserialize l error")
	}
	tokenEndpoints := make([]*TokenTemplateEndPoint, l)
	for i := 0; i < int(l); i++ {
		tt := &TokenTemplateEndPoint{}
		err := tt.Deserialize(source)
		if err != nil {
			return err
		}
		tokenEndpoints[i] = tt
	}
	this.TokenEndpoints = tokenEndpoints
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
	return nil
}

func (this *ResourceDDO) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}

type Fee struct {
	ContractAddr common.Address
	ContractType byte
	Count        uint64
}

func (this *Fee) Serialize(sink *common.ZeroCopySink) {
	sink.WriteAddress(this.ContractAddr)
	sink.WriteByte(this.ContractType)
	sink.WriteUint64(this.Count)
}

type DTokenItem struct {
	Fee         Fee
	ExpiredDate uint64
	Stocks      uint32
	Templates   []TokenTemplate
}

func (this *DTokenItem) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}

func (this *DTokenItem) Serialize(sink *common.ZeroCopySink) {
	this.Fee.Serialize(sink)
	sink.WriteUint64(this.ExpiredDate)
	sink.WriteUint32(this.Stocks)
	sink.WriteVarUint(uint64(len(this.Templates)))
	for _, item := range this.Templates {
		item.Serialize(sink)
	}
}
