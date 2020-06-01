package service

import (
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sellerconfig"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sql"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestMain(t *testing.T) {

	Wallet := sdk.NewWallet("test")
	ServerAccount, err := Wallet.NewDefaultSettingAccount([]byte("123456"))
	assert.Nil(t, err)
	sellerconfig.DefSellerConfig.Wallet = Wallet
	sellerconfig.DefSellerConfig.ServerAccount = ServerAccount
	sellerconfig.DefSellerConfig.Pswd = "123456"
}

func TestSaveDataMeta(t *testing.T) {
	sellerImpl := InitSellerImpl()
	ontId := "did:ont:AcVBV1zKGogf9Q54p1Ve78NSQVU5ZUUGkn"
	DataMeta := map[string]interface{}{
		"1": "first",
		"2": "second",
	}
	h, err := ddxf.HashObject(DataMeta)
	assert.Nil(t, err)
	input := io.SellerSaveDataMetaInput{
		DataMeta:     DataMeta,
		DataMetaHash: h,
	}

	output := sellerImpl.SaveDataMeta(input, ontId)
	assert.Equal(t, 0, output.Code)
	fmt.Printf("data %s\n", output.Msg)

	dataStore := &io.SellerSaveDataMeta{}
	filter := bson.M{"dataMetaHash": input.DataMetaHash}
	err = sql.FindElt(sql.DataMetaCollection, filter, dataStore)
	assert.Nil(t, err)

	TokenMeta := map[string]interface{}{
		"1": "first",
		"2": "second",
	}
	ht, err := ddxf.HashObject(TokenMeta)
	assert.Nil(t, err)
	inputt := io.SellerSaveTokenMetaInput{
		TokenMeta:     TokenMeta,
		DataMetaHash:  h,
		TokenMetaHash: ht,
	}
	outputt := sellerImpl.SaveTokenMeta(inputt, ontId)
	assert.Equal(t, 0, outputt.Code)
	fmt.Printf("token %s\n", outputt.Msg)

	tokenStore := &io.SellerSaveTokenMeta{}
	filterT := bson.M{"tokenMetaHash": inputt.TokenMetaHash}
	err = sql.FindElt(sql.TokenMetaCollection, filterT, tokenStore)
	assert.Nil(t, err)

	PublishMeta := map[string]interface{}{
		"3": "three",
		"4": "four",
	}

	inputPub := io.SellerPublishMPItemMetaInput{
		ItemMeta:      PublishMeta,
		TokenMetaHash: inputt.TokenMetaHash,
		DataMetaHash:  inputt.DataMetaHash,
		MPEndpoint:    "xxxMPEndpoint",
	}
	resp, err := sellerImpl.PublishMPItemMeta(inputPub, ontId)
	assert.Nil(t, err)
	fmt.Printf("token %s\n", outputt.Msg)

	qrCodes := &qrCode.QrCode{}
	filterQ := bson.M{"qrCodeId": resp.Id}
	err = sql.FindElt(sql.SellerQrCodeCollection, filterQ, qrCodes)
	assert.Nil(t, err)
}
