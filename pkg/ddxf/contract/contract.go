package contract

import (
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
)

func ConstructPublishParam(seller common.Address, template param.TokenTemplate,tokenResourceType []*param.TokenResourceTyEndpoint,
	itemMetaHash common.Uint256, fee param.Fee, expiredDate uint64, stock uint32, resourceId string) ([]byte, []byte, []byte) {
	ddo := param.ResourceDDO{
		TokenResourceTyEndpoints: tokenResourceType,    // RT for tokens
		Manager:                  seller,               // data owner id
		ItemMetaHash:             itemMetaHash,         // required if len(Templates) > 1
		DTC:                      common.ADDRESS_EMPTY, // can be empty
		MP:                       common.ADDRESS_EMPTY, // can be empty
		Split:                    common.ADDRESS_EMPTY,
	}

	item := param.DTokenItem{
		Fee:         fee,
		ExpiredDate: expiredDate,
		Stocks:      stock,
		Templates:   []param.TokenTemplate{template},
	}

	return []byte(resourceId), ddo.ToBytes(), item.ToBytes()
}
