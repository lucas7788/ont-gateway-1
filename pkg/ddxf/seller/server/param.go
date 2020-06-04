package server

import "github.com/ontio/ontology/common"

type ItemMeta struct {
	ItemMetaHash string                 `bson:"itemMetaHash" json:"itemMetaHash"`
	ItemMetaData map[string]interface{} `bson:"itemMetaData" json:"itemMetaData"`
}

type DataIdInfo struct {
	DataId       string
	DataType     byte
	DataMetaHash common.Uint256
	DataHash     common.Uint256
}

func (this DataIdInfo) Serialize(sink *common.ZeroCopySink) {
	sink.WriteString(this.DataId)
	sink.WriteByte(this.DataType)
	sink.WriteHash(this.DataMetaHash)
	sink.WriteHash(this.DataHash)
}

func (this *DataIdInfo) ToBytes() []byte {
	sink := common.NewZeroCopySink(nil)
	this.Serialize(sink)
	return sink.Bytes()
}
