package server

import (
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/core/types"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/key_manager"
)

func GetAccount(userId string) *ontology_go_sdk.Account {
	plainSeed := []byte(defPlainSeed + userId)
	pri, pubKey := key_manager.GetKeyPair(plainSeed)
	addr := types.AddressFromPubKey(pubKey)
	return &ontology_go_sdk.Account{
		PrivateKey: pri,
		PublicKey:  pubKey,
		Address:    addr,
		SigScheme:  signature.SHA256withECDSA,
	}
}
