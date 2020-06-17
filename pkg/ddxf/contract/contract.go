package contract

import (
	"github.com/ont-bizsuite/ddxf-sdk/ddxf_contract"
	"github.com/ontio/ontology/common"
)

func ConstructPublishParam(seller common.Address, template *ddxf_contract.TokenTemplate,
	tokenResourceType []*ddxf_contract.TokenResourceTyEndpoint,
	itemMetaHash common.Uint256, fee ddxf_contract.Fee, expiredDate uint64, stock uint32) ([]byte, []byte) {
	ddo := ddxf_contract.ResourceDDO{
		TokenResourceTyEndpoints: tokenResourceType,    // RT for tokens
		Manager:                  seller,               // data owner id
		ItemMetaHash:             itemMetaHash,         // required if len(Templates) > 1
		DTC:                      common.ADDRESS_EMPTY, // can be empty
		MP:                       common.ADDRESS_EMPTY, // can be empty
		Split:                    common.ADDRESS_EMPTY,
	}

	item := ddxf_contract.DTokenItem{
		Fee:         fee,
		ExpiredDate: expiredDate,
		Stocks:      stock,
		Templates:   []*ddxf_contract.TokenTemplate{template},
	}

	return ddo.ToBytes(), item.ToBytes()
}
