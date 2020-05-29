package server

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

var mpAcc *ontology_go_sdk.Account

func TestMain(m *testing.M) {
	pri, _ := hex.DecodeString("c19f16785b8f3543bbaf5e1dbb5d398dfa6c85aaad54fc9d71203ce83e505c07")
	mpAcc, _ = ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	fmt.Println(mpAcc.Address.ToBase58())
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
	output := AddRegistryService(in)
	assert.NotEqual(t, output.Code, http.StatusInternalServerError)

	sig, _ := signature2.Sign(user, []byte(mpStr))
	rm := io.MPRemoveRegistryInput{
		MP:   mpStr,
		Sign: sig,
	}
	output2 := RemoveRegistryService(rm)
	assert.NotEqual(t, output2.Code, http.StatusInternalServerError)
}

func TestEndpointImpl_GetAuditRule(t *testing.T) {

	in := io.MPEndpointGetFeeInput{}
	output := GetFeeService(in)
	assert.Equal(t, output.Type, io.ONG)

	itemMeta := make(map[string]interface{})
	itemMeta["key"] = "value"

	publishIn := io.MPEndpointPublishItemMetaInput{
		ItemMeta: io.PublishItemMeta{
			ItemMeta:      itemMeta,
			OnchainItemID: "resource_id",
		},
		SignedDDXFTx: "00d2be2ece5ef40100000000000000c2eb0b000000005c645c529cbab407589537ef4b87b84374f23c38f7ca3703030762a2ad4c2f034d0ce8565641b98407e21364746f6b656e53656c6c65725075626c697368067265736f5f377900010000000020000000000000000000000000000000000000000000000000000000000000000000fbe02b027e61a6d7602f26cfa9487fa58ef9ee7208656e64706f696e74010000000020000000000000000000000000000000000000000000000000000000000000000009656e64706f696e7432000000004c00000000000000000000000000000000000000000164000000000000005ebbce5e0000000001000000010020000000000000000000000000000000000000000000000000000000000000000000024140457a6eaf3f2919416f1ffffdc8372003e5b6d3cb155e60671aba6a64cb22ff56c67c96d9dfab28fc0bed7c7f609ccd28307a347ecd0f0818211c41f53cdd16472321025af6199b152051fb7d508d11897f8e95fa4c95aa76f764dda347f59e9db82955ac4140825972e86b25cdb9b7a16c5baf409ec0c9046831be5df304f1de9839c741019deb7bea6721b636eb95a99d9c5799422c5b832d4f8ddff1abaa64035aa2cf59c023210240cf95b7738a102a554f83f6202fb00ed69f4354a69bb832d8df0938512adde9ac",
	}
	output2 := PublishItemMetaService(publishIn)
	fmt.Println(output2)
	assert.NotEqual(t, output2.Code, http.StatusInternalServerError)
}

func TestEndpointImpl_QueryItemMetas(t *testing.T) {
	//timeout := config.Load().MongoConfig.Timeout
	//ctx, cancel := context.WithTimeout(context.Background(), timeout)
	//defer cancel()
	//item := make(map[string]interface{})
	//item["key"] = "val"
	//itemMeta := io.PublishItemMetaHandler{
	//	OnchainItemID:"resource_id",
	//	ItemMeta:item,
	//}
	//
	//_, err := instance.MongoOfficial().Collection(endpointCollectionName).InsertOne(ctx, itemMeta)
	//assert.Nil(t, err)

	in2 := io.MPEndpointQueryItemMetasInput{
		PageNum:  1,
		PageSize: 10,
	}
	output3 := QueryItemMetasService(in2)
	assert.Equal(t, output3.Msg, "")
	idStr := output3.ItemMetas[0].Id.Hex()
	output := GetItemMetaService(io.MPEndpointGetItemMetaInput{
		ItemMetaID: idStr,
	})
	fmt.Println(output.ItemMeta)
	fmt.Println(output3.ItemMetas[0].Item)
	assert.Equal(t, output.ItemMeta.ItemMeta, output3.ItemMetas[0].Item)
}
