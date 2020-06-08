package seller_buyer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func SaveDataMeta() error {
	//ontId := "did:ont:AcVBV1zKGogf9Q54p1Ve78NSQVU5ZUUGkn"
	DataMeta := map[string]interface{}{
		"1": "first5",
		"2": "second4",
	}
	DataMeta["ISDN"] = "hello"
	h, err := ddxf.HashObject(DataMeta)
	if err != nil {
		return err
	}
	dataMetaHash := hex.EncodeToString(h[:])
	fmt.Println("dataMetaHash:", dataMetaHash)
	input := io.SellerSaveDataMetaInput{
		DataMeta:     DataMeta,
		DataMetaHash: dataMetaHash,
		ExpiredDate:  uint64(time.Now().Unix() + 24*60*60),
		Stock:        1000,
	}

	//send req to seller
	data, err := SendPOST(config.SellerUrl+server.SaveDataMetaUrl, input)
	if err != nil {
		return err
	}
	out := io.SellerSaveDataMetaOutput{}
	err = json.Unmarshal(data, &out)
	if err != nil {
		return err
	}
	fmt.Println("DataId:", out.DataId)
	return nil
}

func SaveTokenMeta(dataMetaHash string) error {
	TokenMeta := map[string]interface{}{
		"1": "first",
		"2": "second",
	}
	ht, err := ddxf.HashObject(TokenMeta)
	if err != nil {
		return err
	}
	tokenMetaHash := hex.EncodeToString(ht[:])
	fmt.Println("tokenMetaHash:", tokenMetaHash)
	input := io.SellerSaveTokenMetaInput{
		TokenMeta:     TokenMeta,
		DataMetaHash:  dataMetaHash,
		TokenMetaHash: tokenMetaHash,
		TokenEndpoint: config.SellerUrl,
	}
	_, err = SendPOST(config.SellerUrl+server.SaveTokenMetaUrl, input)
	return err
}

func PublishMeta1(tokenMetaHash string, dataMetaHash string) error {
	PublishMeta := map[string]interface{}{
		"3": "three",
		"4": "four",
	}

	inputPub := io.SellerPublishMPItemMetaInput{
		ItemMeta:      PublishMeta,
		TokenMetaHash: tokenMetaHash,
		DataMetaHash:  dataMetaHash,
		MPEndpoint:    config.PublishItemMetaUrl,
	}
	data, err := SendPOST(config.SellerUrl+server.PublishItemMetaUrl, inputPub)
	if err != nil {
		return err
	}
	res := qrCode.QrCodeResponse{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	fmt.Println("qrCodeId: ", res.Id)
	return nil
}

func PublishMeta() error {
	qrCodeId := "seller_publishab49ad75-d6f6-4687-90e0-d7bc52eedaec"
	qc := qrCode.QrCode{}
	filterD := bson.M{"qrCodeId": qrCodeId}
	err := server.FindElt(server.SellerQrCodeCollection, filterD, &qc)
	if err != nil {
		return err
	}
	resourceIdBytes, ddo, item, err := server.ParsePublishParamFromQrCodeData(qc.QrCodeData)
	if err != nil {
		return err
	}
	resourceId := string(resourceIdBytes)
	fmt.Println("resourceId: ", resourceId)
	tx, err := instance.OntSdk().DDXFContract(2000000, 500,
		nil).BuildTx(server.ServerAccount, "dtokenSellerPublish", []interface{}{resourceIdBytes, ddo, item})
	if err != nil {
		return err
	}
	txHash := tx.Hash()
	fmt.Println("txHash: ", txHash.ToHexString())
	rawTx, _ := tx.IntoImmutable()
	sink := common.NewZeroCopySink(nil)
	rawTx.Serialization(sink)

	pp := io.PublishParam{}
	filterD = bson.M{"qrCodeId": qrCodeId}
	err = server.FindElt(server.PublishParamCollection, filterD, &pp)
	if err != nil {
		return err
	}

	input := io.MPEndpointPublishItemMetaInput{
		SignedDDXFTx: hex.EncodeToString(sink.Bytes()),
		ItemMeta: io.PublishItemMeta{
			OnchainItemID: hex.EncodeToString(resourceIdBytes),
			ItemMeta:      map[string]interface{}{},
		},
		MPEndpoint: pp.Input.MPEndpoint,
	}
	_, err = SendPOST(config.SellerUrl+server.PublishMPItemMetaUrl, input)
	return err
}
