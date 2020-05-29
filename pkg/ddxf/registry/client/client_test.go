package client

import (
	"encoding/hex"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/core/signature"
	"github.com/stretchr/testify/assert"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"testing"
)

func TestSdk(t *testing.T) {
	user := account.NewAccount("")
	mpStr := "mp4"
	pubkey := hex.EncodeToString(keypair.SerializePublicKey(user.PublicKey))
	input := io.RegistryAddEndpointInput{
		MP:       mpStr,
		Endpoint: "endpoint",
		PubKey:   pubkey,
	}
	output := Sdk().AddEndpoint(input)
	if output.Code != 0 {
		fmt.Println(output.Msg)
		return
	}
	assert.Equal(t, output.Code, 0)
	sig, _ := signature.Sign(user, []byte(mpStr))
	rm := io.RegistryRemoveEndpointInput{
		MP:   mpStr,
		Sign: sig,
	}
	output2 := Sdk().RemoveEndpoint(rm)
	assert.Equal(t, output2.Code, 0)
}

func TestRegistryImplClient_QueryEndpoints(t *testing.T) {
	output3 := Sdk().QueryEndpoints(io.RegistryQueryEndpointsInput{})
	if output3.Code != 0 {
		fmt.Println(output3.Msg)
		return
	}
	assert.Equal(t, output3.Code, 0)
}
