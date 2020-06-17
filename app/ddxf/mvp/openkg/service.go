package main

import (
	"github.com/ont-bizsuite/ddxf-sdk/ddxf_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/core/types"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/key_manager"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
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
