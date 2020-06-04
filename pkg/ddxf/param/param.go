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

type TokenResourceTyEndpoint struct {
	TokenTemplate TokenTemplate
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
}

func (this *ResourceDDO) Serialize(sink *common.ZeroCopySink) {
	sink.WriteAddress(this.Manager)
	sink.WriteVarUint(uint64(len(this.TokenResourceTyEndpoints)))
	for _, v := range this.TokenResourceTyEndpoints {
		v.Serialize(sink)
	}
	//TODO
	sink.WriteBool(true)
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
