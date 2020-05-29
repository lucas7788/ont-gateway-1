package define

import (
	"github.com/ontio/ontology/common"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"io"
)

type TokenTemplate struct {
	DataIDs   string // can be empty
	TokenHash string
}

func ReadString(source *common.ZeroCopySource) (string, error) {
	data, _, irregular, eof := source.NextString()
	if irregular {
		return "", common.ErrIrregularData
	}
	if eof {
		return "", io.ErrUnexpectedEOF
	}
	return data, nil
}
func ConstructTokensAndEndpoint(data []byte, buyer common.Address, onchainItemId string) ([]io2.EndpointToken, error) {
	source := common.NewZeroCopySource(data)
	l, _, irregular, eof := source.NextVarUint()
	if irregular {
		return nil, common.ErrIrregularData
	}
	if eof {
		return nil, io.ErrUnexpectedEOF
	}
	res := make([]io2.EndpointToken, l)
	for i := 0; i < int(l); i++ {
		tt := &TokenTemplate{}
		err := tt.Deserialize(source)
		if err != nil {
			return nil, err
		}
		endpoint, err := ReadString(source)
		if err != nil {
			return nil, err
		}
		res[i] = io2.EndpointToken{
			Token: io2.Token{
				TokenTemplate: *tt,
				Buyer:         buyer,
				OnchainItemId: onchainItemId,
			},
			Endpoint: endpoint,
		}
	}
	return res, nil
}

func (this *TokenTemplate) Serialize(sink *common.ZeroCopySink) {
	if len(this.DataIDs) == 0 {
		sink.WriteBool(false)
	} else {
		sink.WriteBool(true)
		sink.WriteString(this.DataIDs)
	}
	sink.WriteString(this.TokenHash)
}

func (this *TokenTemplate) Deserialize(source *common.ZeroCopySource) error {
	data, irregular, eof := source.NextBool()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	var err error
	if data {
		this.DataIDs, err = ReadString(source)
		if err != nil {
			return err
		}
	}
	this.TokenHash, err = ReadString(source)
	if err != nil {
		return err
	}
	return nil
}
