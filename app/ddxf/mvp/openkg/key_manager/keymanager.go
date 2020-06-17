package key_manager

import (
	"github.com/ontio/go-bip32"
)

//https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
//https://github.com/satoshilabs/slips/blob/master/slip-0044.md
//https://github.com/FactomProject/FactomDocs/blob/master/wallet_info/wallet_test_vectors.md

const (
	Purpose       uint32 = 0x8000002C
	DefaultChange uint32 = 0x0
	TypeBitcoin   uint32 = 0x80000000
	Account       uint32 = 0x80000000
)

// NewKeyFromMasterKey
func NewKeyFromMasterKey(masterKey *bip32.Key, index uint32) (*bip32.Key, error) {
	child, err := masterKey.NewChildKey(Purpose)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(TypeBitcoin)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(Account)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(DefaultChange)
	if err != nil {
		return nil, err
	}

	child, err = child.NewChildKey(index)
	if err != nil {
		return nil, err
	}

	return child, nil
}
