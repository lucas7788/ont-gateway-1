package main

import (
	"github.com/ont-bizsuite/ddxf-sdk/ddxf_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	common2 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/key_manager"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
)

func PublishService(input PublishInput) (output PublishOutput, err error) {
	output.ReqID = input.ReqID
	defer callback(output)

	// 抽取openKGID
	openKGID := input.OpenKGID

	filter := bson.M{"open_kg_id": openKGID}
	param := PublishInput{}
	err = FindElt(OpenKgPublishParamCollection, filter, &param)
	if err != nil {
		return
	}
	// 1. 抽取data meta
	dataMetas := input.Datas

	plainSeed := []byte(defPlainSeed + input.UserID)
	pri, _ := key_manager.GetSerializedKeyPair(plainSeed)

	var seller *ontology_go_sdk.Account

	seller, err = ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	if err != nil {
		return
	}
	resourceId := common.GenerateUUId(config.UUID_RESOURCE_ID)
	// 2. save data metas and publish item
	// send request to seller
	ddo := ddxf_contract.ResourceDDO{
		Manager:                  seller.Address,
		TokenResourceTyEndpoints: []*ddxf_contract.TokenResourceTyEndpoint{},
	}
	item := ddxf_contract.DTokenItem{}
	split := split_policy_contract.SplitPolicyRegisterParam{}
	var tx *types.MutableTransaction
	if input.Delete {
		tx, err = instance.DDXFSdk().DefDDXFKit().BuildFreezePublishTx(common2.ADDRESS_EMPTY,
			[]byte(param.OnChainId), []byte(resourceId), ddo, item, split)
		if err != nil {
			return
		}
		tx, err = instance.DDXFSdk().SignTx(tx, seller)
		if err != nil {
			return
		}
	} else {
		tx, err = instance.DDXFSdk().DefDDXFKit().BuildPublishTx([]byte(resourceId), ddo, item, split)
		if err != nil {
			return
		}
		tx, err = instance.DDXFSdk().SignTx(tx, seller)
		if err != nil {
			return
		}
	}

	// send req to seller
	_, _, _, err = forward.PostJSONRequest("", []byte{}, nil)
	if err != nil {
		return
	}
	return
}

func buyAndUseService() {

}
