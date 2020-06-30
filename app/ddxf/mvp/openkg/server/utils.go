package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/ontid"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/config"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/key_manager"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"net/http"
)

func GetAccount(userId string) *ontology_go_sdk.Account {
	plainSeed := []byte(defPlainSeed + userId)
	pri, pubKey := key_manager.GetKeyPair(plainSeed)
	addr := types.AddressFromPubKey(pubKey)
	return &ontology_go_sdk.Account{
		PrivateKey: pri,
		PublicKey:  pubKey,
		Address:    addr,
		SigScheme:  signature.SHA256withECDSA,
	}
}

func buildAddAttributeTx(hash,dataHash common.Uint256, ontID string,
	seller *ontology_go_sdk.Account)(iMutTx *types.Transaction, err error) {
	var (
		txMut  *types.MutableTransaction
	)

	attr := &ontology_go_sdk.DDOAttribute{
		Key:       []byte("DataMetaHash"),
		Value:     hash[:],
		ValueType: []byte{},
	}
	attr2 := &ontology_go_sdk.DDOAttribute{
		Key:       []byte("DataHash"),
		Value:     dataHash[:],
		ValueType: []byte{},
	}
	attrs := make([]*ontology_go_sdk.DDOAttribute, 0)
	attrs = append(attrs, attr)
	attrs = append(attrs, attr2)

	txMut, err = instance.OntSdk().GetKit().Native.OntId.NewAddAttributesTransaction(
		500, 2000000, ontID, attrs, seller.PublicKey,
	)
	if err != nil {
		return
	}
	err = instance.DDXFSdk().SignTx(txMut, seller)
	if err != nil {
		return
	}
	iMutTx, err = txMut.IntoImmutable()
	if err != nil {
		return
	}
	txHash := txMut.Hash()
	fmt.Println("txhash:", txHash.ToHexString())
	return
}
func regIdWithController(dataId string,seller *ontology_go_sdk.Account) (err error) {
	g := &ontid.Group{
		Members:   []interface{}{},
		Threshold: 1,
	}
	controller := config.PreOntId + seller.Address.ToBase58()
	signers := []ontid.Signer{ontid.Signer{
		Id:    []byte(controller),
		Index: 1,
	}}
	tx, err := instance.DDXFSdk().GetOntologySdk().Native.OntId.NewRegIDWithControllerTransaction(config.GasPrice,
		config.GasLimit, dataId, g, signers)
	if err != nil {
		return
	}
	err = instance.DDXFSdk().SignTx(tx, seller)
	if err != nil {
		return
	}
	imutTx, err := tx.IntoImmutable()
	if err != nil {
		return
	}
	//send tx to seller
	ri := server.RegisterOntIdInput{
		SignedTx: hex.EncodeToString(common.SerializeToBytes(imutTx)),
	}
	data, err := json.Marshal(ri)
	if err != nil {
		return
	}
	code, _, _, err := forward.PostJSONRequest(config.SellerUrl+server.RegisterOntId, data, nil)
	if err != nil {
		return
	}
	txHash := tx.Hash()
	if code != http.StatusOK {
		err = fmt.Errorf("register ontid tx failed, txHash: %s", txHash.ToHexString())
		return
	}
	return
}

func deletePublish(resourceId string, seller *ontology_go_sdk.Account, headers map[string]string) (err error) {
	var (
		tx       *types.MutableTransaction
		iMutTx   *types.Transaction
		bs, data []byte
	)
	tx, err = instance.DDXFSdk().DefMpKit().BuildDeleteTx([]byte(resourceId))
	if err != nil {
		return
	}
	err = instance.DDXFSdk().SignTx(tx, seller)
	if err != nil {
		return
	}

	iMutTx, err = tx.IntoImmutable()
	if err != nil {
		return
	}

	param := server.DeleteParam{SignedTx: hex.EncodeToString(common.SerializeToBytes(iMutTx))}
	bs, err = json.Marshal(param)
	if err != nil {
		return
	}
	_, _, data, err = forward.PostJSONRequestWithRetry(config.SellerUrl+server.FreezeUrl, bs, headers, 10)
	if err != nil {
		return
	}

	res := server.DeleteOutput{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	err = res.Error()
}

func queryDataIdFromSeller(dataMetas []map[string]interface{}) (map[string]interface{},[][sha256.Size]byte, error) {
	dataMetaHashArray := make([]string, len(dataMetas))
	dataMetaHashArray2 := make([][sha256.Size]byte, len(dataMetas))
	for i := 0; i < len(dataMetas); i++ {
		if dataMetas[i]["url"] == nil {
			return nil,nil, errors.New("url empty")
		}
		var hash [sha256.Size]byte
		hash, err := ddxf.HashObject(dataMetas[i])
		if err != nil {
			return nil,nil, err
		}
		dataMetaHashArray[i] = string(hash[:])
		dataMetaHashArray2[i] = hash
	}
	getDataIdParam := server.GetDataIdParam{
		DataMetaHashArray: dataMetaHashArray,
	}
	var data, paramBs []byte
	paramBs, err := json.Marshal(getDataIdParam)
	if err != nil {
		return nil,nil, err
	}
	//查询哪些data id需要上链
	_, _, data, err = forward.PostJSONRequestWithRetry(config.SellerUrl+server.GetDataIdByDataMetaHashUrl, paramBs, headers, 10)
	if err != nil {
		return nil,nil, err
	}
	res := make(map[string]interface{})
	if data != nil {
		err = json.Unmarshal(data, &res)
		if err != nil {
			return nil,nil, err
		}
	}
	return res,dataMetaHashArray2, nil
}
