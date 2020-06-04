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

// ResourceDDO is ddo for resource
type ResourceDDO struct {
	ResourceType      byte
	TokenResourceType map[TokenTemplate]byte   // RT for tokens
	Manager           common.Address           // data owner id
	Endpoint          string                   // data service provider uri
	TokenEndpoint     map[TokenTemplate]string // endpoint for tokens
	DescHash          common.Uint256           // required if len(Templates) > 1
	DTC               common.Address           // can be empty
	MP                common.Address           // can be empty
	Split             common.Address           // can be empty
}

func (this *ResourceDDO) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}
func (this *ResourceDDO) Deserialize(source *common.ZeroCopySource) error {
	source.NextByte()
	l, eof := source.NextUint32()
	if eof {
		return errors.New("1. ResourceDDO Deserialize l error")
	}
	for i := 0; i < int(l); i++ {
		tt := TokenTemplate{}
		tt.Deserialize(source)
		source.NextByte()
	}
	source.NextAddress()
	source.NextString()
	l, eof = source.NextUint32()
	if eof {
		return errors.New("2. ResourceDDO Deserialize l error")
	}
	for i := 0; i < int(l); i++ {
		tt := TokenTemplate{}
		tt.Deserialize(source)
		source.NextString()
	}
	source.NextString()
	var irregular bool
	this.DescHash, eof = source.NextHash()
	if irregular || eof {
		return errors.New("2. ResourceDDO Deserialize l error")
	}
	return nil
}

func (this *ResourceDDO) Serialize(sink *common.ZeroCopySink) {
	sink.WriteByte(0)
	sink.WriteUint32(uint32(len(this.TokenResourceType)))
	for k, v := range this.TokenResourceType {
		k.Serialize(sink)
		sink.WriteByte(v)
	}
	sink.WriteAddress(this.Manager)
	sink.WriteString(this.Endpoint)
	sink.WriteUint32(uint32(len(this.TokenEndpoint)))
	for k, v := range this.TokenEndpoint {
		k.Serialize(sink)
		sink.WriteString(v)
	}
	//TODO
	sink.WriteBool(true)
	sink.WriteHash(this.DescHash)
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
