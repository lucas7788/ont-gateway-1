package mp

import (
	"encoding/hex"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/account"
	signature2 "github.com/ontio/ontology/core/signature"
	"github.com/stretchr/testify/assert"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"net/http"
	"testing"
)

var (
	mp *MarketplaceImpl
)

func TestMain(m *testing.M) {
	pri, _ := hex.DecodeString("c19f16785b8f3543bbaf5e1dbb5d398dfa6c85aaad54fc9d71203ce83e505c07")
	acc, _ := ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	fmt.Println(acc.Address.ToBase58())
	mp = NewMarketplaceImpl(acc)
	mp.Init()
	m.Run()
}

func TestMarketplaceImpl_AddRegistry(t *testing.T) {
	user := account.NewAccount("")
	mpStr := "mp"
	pubkey := hex.EncodeToString(keypair.SerializePublicKey(user.PublicKey))
	in := io.MPAddRegistryInput{
		MP:       mpStr,
		Endpoint: "endpoint",
		PubKey:   pubkey,
	}
	output := mp.AddRegistry(in)
	assert.NotEqual(t, output.Code, http.StatusInternalServerError)

	sig, _ := signature2.Sign(user, []byte(mpStr))
	rm := io.MPRemoveRegistryInput{
		MP:   mpStr,
		Sign: sig,
	}
	output2 := mp.RemoveRegistry(rm)
	assert.NotEqual(t, output2.Code, http.StatusInternalServerError)
}

func TestEndpointImpl_GetAuditRule(t *testing.T) {
	en := mp.Endpoint()

	in := io.MPEndpointGetFeeInput{}
	output := en.GetFee(in)
	assert.Equal(t, output.Type, io.ONG)

	itemMeta := make(map[string]interface{})
	itemMeta["key"] = "value"

	publishIn := io.MPEndpointPublishItemMetaInput{
		ItemMeta:     itemMeta,
		SignedDDXFTx: "",
	}
	en.PublishItemMeta(publishIn)
}
