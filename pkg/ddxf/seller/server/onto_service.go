package server

import (
	"errors"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/contract"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/utils"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func PublishMetaService(input io.SellerPublishMPItemMetaInput, ontId string) (qrCode.QrCodeResponse, error) {
	adT := &io.SellerSaveTokenMeta{}
	filterT := bson.M{"tokenMetaHash": input.TokenMetaHash, "ontId": ontId}
	err := FindElt(TokenMetaCollection, filterT, adT)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}

	adD := &io.SellerSaveDataMeta{}
	filterD := bson.M{"dataMetaHash": input.DataMetaHash, "ontId": ontId}
	err = FindElt(DataMetaCollection, filterD, adD)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}

	arr := strings.Split(ontId, ":")
	if len(arr) != 3 {
		return qrCode.QrCodeResponse{}, err
	}
	sellerAddress, err := common.AddressFromBase58(arr[2])

	// dataMeta related in data contract tx.
	tokenTemplate := &param.TokenTemplate{
		DataID:     adD.DataIds,
		TokenHashs: []string{adT.TokenMetaHash},
	}
	bs, err := ddxf.HashObject(input.ItemMeta)
	itemMetaHash, err := common.Uint256ParseFromBytes(bs[:])
	im := ItemMeta{
		ItemMetaHash: itemMetaHash.ToHexString(),
		ItemMetaData: input.ItemMeta,
	}
	err = InsertElt(ItemMetaCollection, im)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	trt := &param.TokenResourceTyEndpoint{
		TokenTemplate: tokenTemplate,
		ResourceType:  adD.ResourceType,
		Endpoint:      adT.TokenEndpoint,
	}
	resourceIdBytes, resourceDDOBytes, itemBytes := contract.ConstructPublishParam(sellerAddress,
		tokenTemplate,
		[]*param.TokenResourceTyEndpoint{trt},
		itemMetaHash, adD.Fee, adD.ExpiredDate, adD.Stock, adD.DataIds)
	//TODO
	var netType string
	if config.Load().Prod {
		netType = "testnet"
	} else {
		netType = "mainnet"
	}
	qrCodex, err := BuildPublishQrCode(netType, input.MPContractHash,
		resourceIdBytes, resourceDDOBytes, itemBytes, arr[2], ontId)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}

	p := io.PublishParam{
		QrCodeId: qrCodex.QrCodeId,
		Input:    input,
	}
	err = InsertElt(PublishParamCollection, p)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	err = InsertElt(SellerQrCodeCollection, qrCodex)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	return BuildQrCodeResponse(qrCodex.QrCodeId), nil
}

func GetQrCodeByQrCodeIdService(qrCodeId string) (qrCode.QrCode, error) {
	filter := bson.M{"qrCodeId": qrCodeId}
	code := qrCode.QrCode{}
	err := FindElt(SellerQrCodeCollection, filter, &code)
	return code, err
}

func QrCodeCallBackService(param qrCode.QrCodeCallBackParam) error {
	filter := bson.M{"qrCodeId": param.ExtraData.Id}
	code := qrCode.QrCode{}
	err := FindElt(SellerQrCodeCollection, filter, &code)
	if err != nil {
		return err
	}
	uuidType := utils.UUIDType(code.QrCodeId)
	switch uuidType {
	case utils.UUID_TOKEN_SELLER_PUBLISH:
		resourceId, ddo,_, err := ParseFromBytes(code.QrCodeData)
		if err != nil {
			return err
		}
		adD := ItemMeta{}
		filterD := bson.M{"dataMetaHash": ddo.ItemMetaHash}
		err = FindElt(ItemMetaCollection, filterD, &adD)
		if err != nil {
			return err
		}
		pp := io.PublishParam{}
		filterD = bson.M{"qrCodeId": param.ExtraData.Id}
		err = FindElt(PublishParamCollection, filterD, &pp)
		if err != nil {
			return err
		}
		in := io.MPEndpointPublishItemMetaInput{
			SignedDDXFTx: param.SignedTx,
			ItemMeta: io.PublishItemMeta{
				OnchainItemID: resourceId,
				ItemMeta:      adD.ItemMetaData,
			},
			MPEndpoint:pp.Input.MPEndpoint,
		}
		output := PublishMPItemMetaService(in, param.ExtraData.OntId)
		return output.Error()
	default:
		return errors.New("wrong uuid type")
	}
}