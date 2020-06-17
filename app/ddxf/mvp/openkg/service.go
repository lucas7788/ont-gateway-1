package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ont-bizsuite/ddxf-sdk/data_id_contract"
	"github.com/ont-bizsuite/ddxf-sdk/ddxf_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	common2 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/key_manager"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func PublishService(input PublishInput) (output PublishOutput) {
	output.ReqID = input.ReqID
	defer callback(output)

	// 抽取openKGID
	openKGID := input.OpenKGID
	filter := bson.M{"open_kg_id": openKGID}
	param := PublishInput{}
	var findError error
	findError = FindElt(OpenKgPublishParamCollection, filter, &param)
	if findError != nil && findError != mongo.ErrNoDocuments {
		output.Code = http.StatusInternalServerError
		output.Msg = findError.Error()
		return
	}
	plainSeed := []byte(defPlainSeed + input.UserID)
	pri, _ := key_manager.GetSerializedKeyPair(plainSeed)
	seller, err := ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	if input.Delete {
		var tx *types.MutableTransaction
		tx, err = instance.DDXFSdk().DefDDXFKit().BuildFreezeTx([]byte(param.OnChainId))
		if err != nil {
			output.Msg = err.Error()
			output.Code = http.StatusInternalServerError
			return
		}
		tx, err = instance.DDXFSdk().SignTx(tx, seller)
		if err != nil {
			output.Msg = err.Error()
			output.Code = http.StatusInternalServerError
			return
		}
	} else {
		// 1. 抽取data meta
		dataMetas := input.Datas
		dataMetaHashArray := make([]string, len(dataMetas))
		for i := 0; i < len(dataMetas); i++ {
			hash, err := ddxf.HashObject(dataMetas[i])
			if err != nil {
				return
			}
			dataMetaHashArray[i] = string(hash[:])
		}
		getDataIdParam := server.GetDataIdParam{
			DataMetaHashArray: dataMetaHashArray,
		}
		paramBs, err := json.Marshal(getDataIdParam)
		if err != nil {
			return
		}
		_, _, data, err := forward.PostJSONRequest(config.SellerUrl+server.GetDataIdByDataMetaHashUrl, paramBs, nil)
		if err != nil {
			return
		}
		res := server.GetDataIdRes{}
		err = json.Unmarshal(data, res)
		if err != nil {
			return
		}

		ones := make([]io.DataMetaOne, 0)
		for i := 0; i < len(dataMetas); i++ {
			dataMetaHash, err := ddxf.HashObject(dataMetas[i])
			if err != nil {
				return
			}
			for _, hash := range res.DataIdAndDataMetaHashArray {
				if hash.DataMetaHash == string(dataMetaHash[:]) {
					if hash.DataId == "" {
						dataId := common.GenerateUUId(config.UUID_PRE_DATAID)
						one := io.DataMetaOne{
							DataMeta:     dataMetas[i],
							DataMetaHash: string(dataMetaHash[:]),
							DataEndpoint: config.SellerUrl,
							ResourceType: 0,
							DataHash:     "",
							DataId:       dataId,
						}
						ones = append(ones, one)
					}
				}
			}
		}
		// invoke seller saveDataMeta
		infos := make([]data_id_contract.DataIdInfo, len(ones))
		for i := 0; i < len(ones); i++ {
			hash, err := common2.Uint256FromHexString(ones[i].DataMetaHash)
			if err != nil {
				return
			}
			dataHash, err := common2.Uint256FromHexString(ones[i].DataHash)
			infos[i] = data_id_contract.DataIdInfo{
				DataId:       ones[i].DataId,
				DataType:     ones[i].ResourceType,
				DataMetaHash: hash,
				DataHash:     dataHash,
				Owners:       []*data_id_contract.OntIdIndex{},
			}
		}
		tx, err := instance.DDXFSdk().DefDataIdKit().BuildRegisterDataIdInfoArrayTx(infos)
		if err != nil {
			return
		}
		tx, err = instance.DDXFSdk().SignTx(tx, seller)
		if err != nil {
			return
		}
		iMutTx, err := tx.IntoImmutable()
		saveDataMetaArray := io.SellerSaveDataMetaArrayInput{
			DataMetaOneArray: ones,
			SignedTx:         hex.EncodeToString(common2.SerializeToBytes(iMutTx)),
		}

		bs, err := json.Marshal(saveDataMetaArray)
		if err != nil {
			return
		}
		_, _, data, err = forward.PostJSONRequest(config.SellerUrl+server.SaveDataMetaArrayUrl, bs, nil)

		trte := make([]*ddxf_contract.TokenResourceTyEndpoint, len(dataMetas))
		for i := 0; i < len(dataMetas); i++ {
			trte[1] = &ddxf_contract.TokenResourceTyEndpoint{
				TokenTemplate: &ddxf_contract.TokenTemplate{
					DataID:     "",
					TokenHashs: []string{},
				},
				ResourceType: 0,
				Endpoint:     config.SellerUrl,
			}
		}

		resourceId := common.GenerateUUId(config.UUID_RESOURCE_ID)
		input.OnChainId = resourceId
		InsertElt(OpenKgPublishParamCollection, input)
		// 2. save data metas and publish item
		// send request to seller
		var itemMetaHash [32]byte
		itemMetaHash, err = ddxf.HashObject(input.Item)
		ddo := ddxf_contract.ResourceDDO{
			Manager:                  seller.Address,
			TokenResourceTyEndpoints: trte,
			ItemMetaHash:             itemMetaHash,
		}

		item := ddxf_contract.DTokenItem{
			Fee:         ddxf_contract.Fee{},
			ExpiredDate: uint64(time.Now().Unix()) + uint64(time.Hour*24*30),
			Stocks:      1000,
			Templates: []*ddxf_contract.TokenTemplate{
				&ddxf_contract.TokenTemplate{
					DataID:     "",
					TokenHashs: []string{},
				},
			},
		}
		split := split_policy_contract.SplitPolicyRegisterParam{}
		var tx *types.MutableTransaction
		if findError == mongo.ErrNilDocument {
			tx, err = instance.DDXFSdk().DefDDXFKit().BuildPublishTx([]byte(resourceId), ddo, item, split)
			if err != nil {
				output.Code = http.StatusInternalServerError
				output.Msg = err.Error()
				return
			}
			tx, err = instance.DDXFSdk().SignTx(tx, seller)
			if err != nil {
				output.Code = http.StatusInternalServerError
				output.Msg = err.Error()
				return
			}
		} else {
			tx, err = instance.DDXFSdk().DefDDXFKit().BuildFreezeAndPublishTx(
				[]byte(param.OnChainId), []byte(resourceId), ddo, item, split)
			if err != nil {
				output.Code = http.StatusInternalServerError
				output.Msg = err.Error()
				return
			}
			tx, err = instance.DDXFSdk().SignTx(tx, seller)
			if err != nil {
				output.Code = http.StatusInternalServerError
				output.Msg = err.Error()
				return
			}
		}
	}

	// send req to seller
	_, _, _, err = forward.PostJSONRequest("", []byte{}, nil)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	return
}

func buyAndUseService(input BuyAndUseInput) (output BuyAndUseOutput) {
	defer callback(output)
	output.ReqID = input.ReqID
	plainSeed := []byte(defPlainSeed + input.UserID)
	pri, _ := key_manager.GetSerializedKeyPair(plainSeed)
	user, err := ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}
	filter := bson.M{"open_kg_id": input.OpenKGID}
	param := PublishInput{}
	err = FindElt(OpenKgPublishParamCollection, filter, &param)
	if param.Delete {
		output.Msg = "has delete"
		output.Code = http.StatusInternalServerError
		return
	}
	tokenTemplate := ddxf_contract.TokenTemplate{
		DataID:     input.DataID,
		TokenHashs: []string{},
	}
	var tx *types.MutableTransaction
	tx, err = instance.DDXFSdk().DefDDXFKit().BuildBuyAndUseTokenTx(user.Address,
		payer.Address, []byte(param.OnChainId), 1, tokenTemplate)
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}
	tx, err = instance.DDXFSdk().SignTx(tx, user)
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}
	tx, err = instance.DDXFSdk().SignTx(tx, payer)
	if err != nil {
		output.Msg = err.Error()
		output.Code = http.StatusInternalServerError
		return
	}
	return
}
