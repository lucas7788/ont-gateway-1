package data_id_contract

import (
	"fmt"
	"github.com/ontio/ontology/common"
)

type DataIdInfo struct {
	DataId       string
	DataMetaHash common.Uint256
	DataHash     common.Uint256
	Owners       []*OntIdIndex
}

type OntIdIndex struct {
	OntId string
	index uint16
}

func (this *OntIdIndex) Serialize(sink *common.ZeroCopySink) {
	sink.WriteString(this.OntId)
	sink.WriteUint16(this.index)
}
func (this *OntIdIndex) Deserialize(source *common.ZeroCopySource) error {
	id, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("read ontid error, irregular :%v,eof: %v", irregular, eof)
	}
	this.OntId = id
	this.index, eof = source.NextUint16()
	if irregular || eof {
		return fmt.Errorf("read ontid index error, eof: %v", eof)
	}
	return nil
}

func (this DataIdInfo) Serialize(sink *common.ZeroCopySink) {
	sink.WriteString(this.DataId)
	sink.WriteHash(this.DataMetaHash)
	sink.WriteHash(this.DataHash)
	sink.WriteVarUint(uint64(len(this.Owners)))
	for _, v := range this.Owners {
		v.Serialize(sink)
	}
}

func (this *DataIdInfo) Deserialize(source *common.ZeroCopySource) error {
	dataId, _, irregular, eof := source.NextString()
	if irregular || eof {
		return fmt.Errorf("read data id error, irregular :%v,eof: %v", irregular, eof)
	}
	this.DataId = dataId
	this.DataMetaHash, eof = source.NextHash()
	this.DataHash, eof = source.NextHash()
	l, _, irregular, eof := source.NextVarUint()
	if irregular || eof {
		return fmt.Errorf("read owner length error, irregular :%v,eof: %v", irregular, eof)
	}
	idx := make([]*OntIdIndex, l)
	for i := 0; i < int(l); i++ {
		err := idx[i].Deserialize(source)
		if err != nil {
			return err
		}
	}
	this.Owners = idx
	return nil
}

func (this *DataIdInfo) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}

func (this *DataIdInfo) FromBytes(data []byte) error {
	source := common.NewZeroCopySource(data)
	return this.Deserialize(source)
}
