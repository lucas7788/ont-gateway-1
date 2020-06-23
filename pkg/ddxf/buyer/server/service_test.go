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
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
)

var (
	resourceId = "resourceid_750773ad-883d-40b7-9473-818c7a578403"
	DataID     = "dataid_1e5d9f5b-9332-4145-8a59-98360716ae01"
	tokenHash  = "e2a740fa12bd94f0e242688e29f6d803f7671eb1f81bcfbdc1c3e213878e7dd4"
)

func TestBuildBuyGetQrCodeRsp(t *testing.T) {
	bs, _ := hex.DecodeString("7265736f7572636569645f64303561383338642d323132332d343835642d616632642d383932396235626562666436")
	fmt.Println(string(bs))
}

func TestBuyDTokenService(t *testing.T) {
	buyParam := []interface{}{resourceId, 1, BuyerMgrAccount.Address}
	tx, err := instance.OntSdk().DDXFContract(2000000, 500,
		nil).BuildTx(BuyerMgrAccount, "BuyDtoken", buyParam)
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
	tokenHashBytes, _ := hex.DecodeString(tokenHash)
	template := market_place_contract.TokenTemplate{
		DataID:     DataID,
		TokenHashs: []string{string(tokenHashBytes)},
	}
	fmt.Println(hex.EncodeToString(template.ToBytes()))
	fmt.Println(BuyerMgrAccount.Address.ToBase58())
	userTokenParam := []interface{}{resourceId, BuyerMgrAccount.Address, template.ToBytes(), 1}
	tx, _ := instance.OntSdk().DDXFContract(2000000,
		500, nil).BuildTx(BuyerMgrAccount, "useToken", userTokenParam)

	txhash := tx.Hash()
	fmt.Println("txhash:", txhash.ToHexString())
	imMut, _ := tx.IntoImmutable()
	txHash := tx.Hash()
	fmt.Println("txHash:", txHash.ToHexString())
	sink := common2.NewZeroCopySink(nil)
	imMut.Serialization(sink)
	input := io.BuyerUseTokenInput{
		Tx:              hex.EncodeToString(sink.Bytes()),
		TokenOpEndpoint: config.SellerUrl,
	}
	output := UseTokenService(input)
	assert.Nil(t, output.Error())
}

func TestHandleEvent(t *testing.T) {
	res, err := common.HandleEvent("0f792177d846c2e4a69e0a7a2058ced610febf701e8a671a9b0cb4447a5e1416", "BuyDtoken")
	assert.Nil(t, err)
	fmt.Println(res)
}
