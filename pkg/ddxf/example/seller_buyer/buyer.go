package seller_buyer

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/buyer/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
)

func BuyDtoken(buyer *ontology_go_sdk.Account, resourceId string) error {
	buyParam := []interface{}{resourceId, 1, buyer.Address}
	tx, err := instance.OntSdk().DDXFContract(2000000, 500,
		nil).BuildTx(buyer, "buyDtoken", buyParam)
	if err != nil {
		return err
	}

	fmt.Println("[BuyDtoken] buyerAddress:", buyer.Address.ToBase58())
	txImu, _ := tx.IntoImmutable()
	sink := common.NewZeroCopySink(nil)
	txImu.Serialization(sink)
	input := io.BuyerBuyDtokenInput{
		SignedTx: hex.EncodeToString(sink.Bytes()),
	}
	txHash := tx.Hash()
	fmt.Println("[BuyDtoken] txHash:", txHash.ToHexString())
	_, err = SendPOST(config.BuyerUrl+server.BuyDtoken, input)
	return err
}

func UseToken(buyer *ontology_go_sdk.Account, resourceId, tokenMetaHash string, dataId string) error {
	tokenHashBytes, _ := hex.DecodeString(tokenMetaHash)
	template := &market_place_contract.TokenTemplate{
		DataID:     dataId,
		TokenHashs: []string{string(tokenHashBytes)},
	}
	fmt.Println("[UseToken] template: ", hex.EncodeToString(template.ToBytes()))
	userTokenParam := []interface{}{resourceId, buyer.Address, template.ToBytes(), 1}
	tx, err := instance.OntSdk().DefaultDDXFContract().BuildTx(buyer, "useToken", userTokenParam)
	if err != nil {
		return err
	}
	txhash := tx.Hash()
	fmt.Println("[UseToken] txhash:", txhash.ToHexString())
	imMut, _ := tx.IntoImmutable()
	sink := common.NewZeroCopySink(nil)
	imMut.Serialization(sink)
	input := io.BuyerUseTokenInput{
		Tx:              hex.EncodeToString(sink.Bytes()),
		TokenOpEndpoint: config.SellerUrl,
	}
	fmt.Println("[UseToken] input: ", input)
	var data []byte
	data, err = SendPOST(config.BuyerUrl+server.UseDToken, input)
	if err != nil {
		return err
	}
	rrr := make(map[string]interface{})
	err = json.Unmarshal(data, &rrr)
	if err != nil {
		return err
	}
	bs, err := base64.RawURLEncoding.DecodeString(rrr["Result"].(string))
	if err != nil {
		return err
	}
	fmt.Println("[UseToken] buyer use token result: ", string(bs))
	return err
}
