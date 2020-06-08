package seller_buyer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ddxf"
	common2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/contract"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"math/rand"
	"strconv"
	"time"
)

func SaveDataMeta() (*io.SellerSaveDataMetaOutput, *io.SellerSaveDataMetaInput, error) {
	//ontId := "did:ont:AcVBV1zKGogf9Q54p1Ve78NSQVU5ZUUGkn"
	i := rand.Int()
	DataMeta := map[string]interface{}{
		"1": "first6" + strconv.Itoa(i),
		"2": "second6",
	}
	DataMeta["ISDN"] = "hello"
	h, err := ddxf.HashObject(DataMeta)
	if err != nil {
		return nil, nil, err
	}
	dataMetaHash := hex.EncodeToString(h[:])
	fmt.Println("dataMetaHash:", dataMetaHash)
	input := &io.SellerSaveDataMetaInput{
		DataMeta:     DataMeta,
		DataMetaHash: dataMetaHash,
		ExpiredDate:  uint64(time.Now().Unix() + 24*60*60),
		Stock:        1000,
	}

	//send req to seller
	data, err := SendPOST(config.SellerUrl+server.SaveDataMetaUrl, input)
	if err != nil {
		return nil, nil, err
	}
	out := &io.SellerSaveDataMetaOutput{}
	err = json.Unmarshal(data, out)
	if err != nil {
		return nil, nil, err
	}
	return out, input, nil
}

func SaveTokenMeta(dataMetaHash string) (input io.SellerSaveTokenMetaInput, err error) {
	TokenMeta := map[string]interface{}{
		"1": "first",
		"2": "second",
	}
	var ht [32]byte
	ht, err = ddxf.HashObject(TokenMeta)
	if err != nil {
		return
	}
	tokenMetaHash := hex.EncodeToString(ht[:])
	fmt.Println("tokenMetaHash:", tokenMetaHash)
	input = io.SellerSaveTokenMetaInput{
		TokenMeta:     TokenMeta,
		DataMetaHash:  dataMetaHash,
		TokenMetaHash: tokenMetaHash,
		TokenEndpoint: config.SellerUrl,
	}
	_, err = SendPOST(config.SellerUrl+server.SaveTokenMetaUrl, input)
	return
}

func PublishMeta(seller *ontology_go_sdk.Account, saveDataMetaOut *io.SellerSaveDataMetaOutput,
	saveDataMetaIn *io.SellerSaveDataMetaInput, saveTokenMetaIn io.SellerSaveTokenMetaInput) (string, error) {

	resourceIdBytes := []byte(common2.GenerateUUId(config.UUID_RESOURCE_ID))
	fmt.Println("resourceId:", string(resourceIdBytes))
	tokenMetaHash, _ := hex.DecodeString(saveTokenMetaIn.TokenMetaHash)
	tokenTemplate := &param.TokenTemplate{
		DataID:     saveDataMetaOut.DataId,
		TokenHashs: []string{string(tokenMetaHash)},
	}
	itemMeta := map[string]interface{}{
		"3": "three",
		"4": "four",
	}

	bs, err := ddxf.HashObject(itemMeta)
	itemMetaHash, err := common.Uint256ParseFromBytes(bs[:])
	resourceDDOBytes, itemBytes := contract.ConstructPublishParam(seller.Address,
		tokenTemplate,
		[]*param.TokenResourceTyEndpoint{&param.TokenResourceTyEndpoint{
			TokenTemplate: tokenTemplate,
			ResourceType:  saveDataMetaIn.ResourceType,
			Endpoint:      saveDataMetaIn.DataEndpoint,
		}},
		itemMetaHash, saveDataMetaIn.Fee, saveDataMetaIn.ExpiredDate, saveDataMetaIn.Stock)
	tx, err := instance.OntSdk().DefaultDDXFContract().BuildTx(seller, "dtokenSellerPublish",
		[]interface{}{resourceIdBytes, resourceDDOBytes, itemBytes})
	if err != nil {
		return "", err
	}
	txHash := tx.Hash()
	fmt.Println("txHash: ", txHash.ToHexString())
	rawTx, _ := tx.IntoImmutable()
	sink := common.NewZeroCopySink(nil)
	rawTx.Serialization(sink)

	input := io.MPEndpointPublishItemMetaInput{
		SignedDDXFTx: hex.EncodeToString(sink.Bytes()),
		ItemMeta: io.PublishItemMeta{
			OnchainItemID: hex.EncodeToString(resourceIdBytes),
			ItemMeta:      map[string]interface{}{},
		},
		MPEndpoint: config.PublishItemMetaUrl,
	}
	_, err = SendPOST(config.SellerUrl+server.PublishMPItemMetaUrl, input)
	return string(resourceIdBytes), err
}
