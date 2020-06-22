package config

import (
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"encoding/hex"
)

var defDDXFConfig *DDXFConfig

type DDXFConfig struct {
	OperatorAccount *ontology_go_sdk.Account
	OperatorOntid   string
}

func DefDDXFConfig() *DDXFConfig {
	if defDDXFConfig == nil {
		pri, _ := hex.DecodeString("c19f16785b8f3543bbaf5e1dbb5d398dfa6c85aaad54fc9d71203ce83e505c07")
		acc, _ := ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
		defDDXFConfig = &DDXFConfig{
			OperatorAccount: acc,
			OperatorOntid:   "did:id:",
		}
	}
	return defDDXFConfig
}
