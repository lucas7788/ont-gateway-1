package contract

import (
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
)

func ConstructPublishParam(seller common.Address, template param.TokenTemplate, tokenEndpointUrl string,
	itemMetaHash common.Uint256, resourceType byte, fee param.Fee, expiredDate uint64, stock uint32, resourceId string) ([]byte, []byte, []byte) {
	tokenResourceType := make(map[param.TokenTemplate]byte)
	tokenResourceType[template] = resourceType
	tokenEndpoint := make(map[param.TokenTemplate]string)
	tokenEndpoint[template] = tokenEndpointUrl
	ddo := param.ResourceDDO{
		ResourceType:      resourceType,
		TokenResourceType: tokenResourceType,    // RT for tokens
		Manager:           seller,               // data owner id
		Endpoint:          tokenEndpointUrl,     // data service provider uri
		TokenEndpoint:     tokenEndpoint,        // endpoint for tokens
		DescHash:          itemMetaHash,         // required if len(Templates) > 1
		DTC:               common.ADDRESS_EMPTY, // can be empty
		MP:                common.ADDRESS_EMPTY, // can be empty
		Split:             common.ADDRESS_EMPTY,
	}

	item := param.DTokenItem{
		Fee:         fee,
		ExpiredDate: expiredDate,
		Stocks:      stock,
		Templates:   []param.TokenTemplate{template},
	}

	return []byte(resourceId), ddo.ToBytes(), item.ToBytes()
}
