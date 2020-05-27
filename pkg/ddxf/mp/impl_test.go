package mp

import (
	"testing"

	"encoding/hex"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/account"
	signature2 "github.com/ontio/ontology/core/signature"
	"github.com/stretchr/testify/assert"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"net/http"
)

var (
	mp *MarketplaceImpl
)

func TestMain(m *testing.M) {
	acc := ontology_go_sdk.NewAccount(signature.SHA256withECDSA)
	mp = NewMarketplaceImpl(acc)
}

func TestMarketplaceImpl_AddRegistry(t *testing.T) {
	user := account.NewAccount("")
	mpStr := "mp"
	in := io.MPAddRegistryInput{
		MP:       mpStr,
		Endpoint: "endpoint",
		PubKey:   hex.EncodeToString(keypair.SerializePublicKey(user.PublicKey)),
	}
	output := mp.AddRegistry(in)
	assert.NotEqual(t, output.Code, http.StatusInternalServerError)

	sig, _ := signature2.Sign(user, []byte(""))

	rm := io.MPRemoveRegistryInput{
		MP:   mpStr,
		Sign: sig,
	}
	output2 := mp.RemoveRegistry(rm)
	assert.NotEqual(t, output2.Code, http.StatusInternalServerError)
}
