package server

import (
	"testing"

	"encoding/hex"
	"fmt"
	common2 "github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
)

var resourceId = "resourceId2"

func TestBuyDTokenService(t *testing.T) {
	buyParam := []interface{}{resourceId, 1, BuyerMgrAccount.Address}
	tx, err := instance.OntSdk().DDXFContract(2000000, 500,
		nil).BuildTx(BuyerMgrAccount, "buyDtoken", buyParam)
	assert.Nil(t, err)

	fmt.Println("buyerAddress:", BuyerMgrAccount.Address.ToBase58())
	txImu, _ := tx.IntoImmutable()
	sink := common2.NewZeroCopySink(nil)
	txImu.Serialization(sink)
	input := io.BuyerBuyDtokenInput{
		SignedTx: hex.EncodeToString(sink.Bytes()),
	}
	txHash := tx.Hash()
	fmt.Println("txHash:", txHash.ToHexString())

	output := BuyDTokenService(input)
	assert.Equal(t, 0, output.Code)
	fmt.Println("BuyDTokenService output: ", output)
}

func TestUseTokenService(t *testing.T) {
	tokenHash := make([]byte, 32)
	template := param.TokenTemplate{
		DataIDs:   "",
		TokenHash: []string{string(tokenHash)},
	}
	userTokenParam := []interface{}{resourceId, BuyerMgrAccount.Address, template.ToBytes(), 1}
	tx, _ := instance.OntSdk().DDXFContract(2000000,
		500, nil).BuildTx(BuyerMgrAccount, "useToken", userTokenParam)

	imMut, _ := tx.IntoImmutable()
	txHash := tx.Hash()
	fmt.Println("txHash:", txHash.ToHexString())
	sink := common2.NewZeroCopySink(nil)
	imMut.Serialization(sink)
	input := io.BuyerUseTokenInput{
		Tx:              hex.EncodeToString(sink.Bytes()),
		TokenOpEndpoint: config.SellerUseTokenUrl,
	}
	output := UseTokenService(input)
	assert.Nil(t, output.Error())
}

func TestHandleEvent(t *testing.T) {
	res, err := common.HandleEvent("0f792177d846c2e4a69e0a7a2058ced610febf701e8a671a9b0cb4447a5e1416", "buyDtoken")
	assert.Nil(t, err)
	fmt.Println(res)
}
