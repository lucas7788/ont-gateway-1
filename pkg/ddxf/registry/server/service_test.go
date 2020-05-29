package server

import (
	"testing"

	"encoding/hex"
	"github.com/magiconair/properties/assert"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/core/signature"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
)

func TestQueryEndpointsService(t *testing.T) {
	user := account.NewAccount("")
	mpStr := "mp2"
	pubkey := hex.EncodeToString(keypair.SerializePublicKey(user.PublicKey))
	input := io.RegistryAddEndpointInput{
		MP:       mpStr,
		Endpoint: "endpoint",
		PubKey:   pubkey,
	}
	output := AddEndpointService(input)
	assert.Equal(t, output.Code, 0)
	output3 := QueryEndpointsService(io.RegistryQueryEndpointsInput{})
	assert.Equal(t, output3.Code, 0)
}

func TestAddEndpointService(t *testing.T) {
	Init()
	user := account.NewAccount("")
	mpStr := "mp"
	pubkey := hex.EncodeToString(keypair.SerializePublicKey(user.PublicKey))
	input := io.RegistryAddEndpointInput{
		MP:       mpStr,
		Endpoint: "endpoint",
		PubKey:   pubkey,
	}
	output := AddEndpointService(input)
	assert.Equal(t, output.Msg, "")
	sig, _ := signature.Sign(user, []byte(mpStr))
	rm := io.RegistryRemoveEndpointInput{
		MP:   mpStr,
		Sign: sig,
	}
	output2 := RemoveEndpointService(rm)
	assert.Equal(t, output2.Msg, "")

}
