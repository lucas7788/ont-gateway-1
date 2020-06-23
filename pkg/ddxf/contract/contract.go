package contract

import (
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ontio/ontology/common"
)

func ConstructPublishParam(seller common.Address, template *market_place_contract.TokenTemplate,
	itemMetaHash common.Uint256, fee market_place_contract.Fee, expiredDate uint64, stock uint32) ([]byte, []byte) {
	ddo := market_place_contract.ResourceDDO{
		Manager:      seller,               // data owner id
		ItemMetaHash: itemMetaHash,         // required if len(Templates) > 1
		DTC:          common.ADDRESS_EMPTY, // can be empty
		MP:           common.ADDRESS_EMPTY, // can be empty
		Split:        common.ADDRESS_EMPTY,
	}

	item := market_place_contract.DTokenItem{
		Fee:         fee,
		ExpiredDate: expiredDate,
		Stocks:      stock,
		Templates:   []*market_place_contract.TokenTemplate{template},
	}

	return ddo.ToBytes(), item.ToBytes()
}
