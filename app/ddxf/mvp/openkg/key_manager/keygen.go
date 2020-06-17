package key_manager

import (
	"crypto/elliptic"
	"fmt"
	"github.com/ontio/go-bip32"
	"github.com/ontio/ontology-crypto/ec"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.uber.org/zap"
)

func GetKeyPair(plainSeed []byte) (privKey keypair.PrivateKey, pubKey keypair.PublicKey) {
	for i := uint32(0); i < 10; i++ {
		pri, pub, err := getKeyPair(plainSeed, i)
		if err != nil {
			instance.Logger().Error("getKeyPair faild:", zap.Error(err))
			i++
			continue
		}
		return pri, pub
	}
	return
}
func GetSerializedKeyPair(seed []byte) (privkey, pubkey []byte) {
	var (
		_privkey keypair.PrivateKey
		_pubkey  keypair.PublicKey
	)
	_privkey, _pubkey = GetKeyPair(seed)

	privkey = keypair.SerializePrivateKey(_privkey)
	pubkey = keypair.SerializePublicKey(_pubkey)

	return
}

func getKeyPair(plainSeed []byte, index uint32) (
	privKey keypair.PrivateKey, pubKey keypair.PublicKey, err error) {

	var (
		keyBytes   []byte
		mk, newkey *bip32.Key
	)
	// defer ClearSeed(plainSeed)

	if mk, err = bip32.NewMasterKey(plainSeed); err != nil {
		fmt.Printf("NewMasterKey error:%s\n", err)
		return
	}
	// defer ClearKey(mk)

	if newkey, err = NewKeyFromMasterKey(mk, index); err != nil {
		fmt.Printf("NewKeyFromMasterKey error:%s\n", err)
		return
	}
	// defer ClearKey(newkey)

	if keyBytes, err = newkey.Serialize(); err != nil {
		fmt.Printf("Serialize error:%s\n", err)
		return
	}
	// defer ClearSeed(keyBytes)

	temp := ec.ConstructPrivateKey(keyBytes[46:78], elliptic.P256())
	privKey = &ec.PrivateKey{
		Algorithm:  ec.ECDSA,
		PrivateKey: temp,
	}
	pubKey = privKey.Public()

	return
}
