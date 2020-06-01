package server

import (
	"testing"

	"encoding/hex"
	"fmt"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	pri, _ := hex.DecodeString("c19f16785b8f3543bbaf5e1dbb5d398dfa6c85aaad54fc9d71203ce83e505c07")
	BuyerMgrAccount, _ = ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	fmt.Println(BuyerMgrAccount.Address.ToBase58())
	m.Run()
}

func TestBuyDtokenQrCodeService(t *testing.T) {
	input := BuyerBuyDtokenQrCodeInput{
		OnchainItemId:   "test",
		N:               1,
		Buyer:           "",
		TokenOpEndpoint: "",
	}
	c, err := BuyDtokenQrCodeService(input)
	assert.Nil(t, err)
	fmt.Println("c: ", c)

	code, err := GetQrCodeByQrCodeIdService(c.Id)
	assert.Nil(t, err)
	fmt.Println("code: ", code)
}

func TestGetQrCodeByQrCodeIdService(t *testing.T) {
	code, err := GetQrCodeByQrCodeIdService("6ea52a5a-ce96-476f-b189-0068fed54430")
	assert.Nil(t, err)
	fmt.Println("code: ", code)
}
