package config

import "github.com/ontio/ontology-go-sdk"

var DefDDXFConfig DDXFConfig

type DDXFConfig struct {
	OperatorAccount *ontology_go_sdk.Account
	OperatorOntid   string
}
