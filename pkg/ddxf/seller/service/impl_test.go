package service

import (
	"encoding/hex"
	"fmt"
	"github.com/ontio/ontology-crypto/signature"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/contract"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sellerconfig"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sql"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestMain(t *testing.M) {
	Wallet := sdk.NewWallet("test")
	var err error
	pri, _ := hex.DecodeString("c19f16785b8f3543bbaf5e1dbb5d398dfa6c85aaad54fc9d71203ce83e505c07")
	ServerAccount, _ = sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	if err != nil {
		panic(err)
	}
	sellerconfig.DefSellerConfig.Wallet = Wallet
	sellerconfig.DefSellerConfig.ServerAccount = ServerAccount
	sellerconfig.DefSellerConfig.Pswd = "123456"
	InitSellerImpl()
	t.Run()
}

func TestSaveDataMeta(t *testing.T) {

	ontId := "did:ont:AcVBV1zKGogf9Q54p1Ve78NSQVU5ZUUGkn"
	DataMeta := map[string]interface{}{
		"1": "first",
		"2": "second",
	}
	h, err := ddxf.HashObject(DataMeta)
	assert.Nil(t, err)
	input := io.SellerSaveDataMetaInput{
		DataMeta:     DataMeta,
		DataMetaHash: hex.EncodeToString(h[:]),
	}

	output := DefSellerImpl.SaveDataMeta(input, ontId)
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
		DataMetaHash:  hex.EncodeToString(h[:]),
		TokenMetaHash: hex.EncodeToString(ht[:]),
	}
	outputt := DefSellerImpl.SaveTokenMeta(inputt, ontId)
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
	res, err := PublishMetaService(inputPub, ontId)
	assert.Nil(t, err)
	fmt.Printf("token %s\n", outputt.Msg)

	qrCodes := &qrCode.QrCode{}
	filterQ := bson.M{"qrCodeId": res.Id}
	err = sql.FindElt(sql.SellerQrCodeCollection, filterQ, qrCodes)
	assert.Nil(t, err)
}

func TestSellerImpl_PublishMPItemMeta(t *testing.T) {
	fmt.Println("seller address:", ServerAccount.Address.ToBase58())

	tokenTemplate := param.TokenTemplate{
		DataIDs:    "",
		TokenHashs: []string{string(common.UINT256_EMPTY[:])},
	}
	itemMetaData := map[string]interface{}{
		"item": "val",
	}

	bs, err := ddxf.HashObject(itemMetaData)
	itemMetaHash, err := common.Uint256ParseFromBytes(bs[:])
	expiredDate := time.Now().Unix() + 10*24*60*60
	resourceId, ddo, item := contract.ConstructPublishParam(ServerAccount.Address, tokenTemplate,
		"tokenendpointurl", itemMetaHash, 0, param.Fee{
			Count: 1,
		}, uint64(expiredDate), 100, "resourceId2")

	tx, err := instance.OntSdk().DDXFContract(2000000, 500,
		nil).BuildTx(ServerAccount, "dtokenSellerPublish", []interface{}{resourceId, ddo, item})
	assert.Nil(t, err)
	rawTx, _ := tx.IntoImmutable()
	sink := common.NewZeroCopySink(nil)
	rawTx.Serialization(sink)
	input := io.MPEndpointPublishItemMetaInput{
		SignedDDXFTx: hex.EncodeToString(sink.Bytes()),
		ItemMeta: io.PublishItemMeta{
			OnchainItemID: hex.EncodeToString(resourceId),
			ItemMeta:      map[string]interface{}{},
		},
	}
	output := DefSellerImpl.PublishMPItemMeta(input, "ontid")
	assert.Nil(t, output.Error())
}
