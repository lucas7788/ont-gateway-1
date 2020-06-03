package service

import (
	"errors"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/contract"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	qrCode2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sellerconfig"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sql"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/utils"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"crypto/sha256"
	"encoding/json"
)

func PublishMetaService(input io.SellerPublishMPItemMetaInput, ontId string) (qrCode.QrCodeResponse, error) {
	adT := &io.SellerSaveTokenMeta{}
	filterT := bson.M{"tokenMetaHash": input.TokenMetaHash, "ontId": ontId}
	err := sql.FindElt(sql.TokenMetaCollection, filterT, adT)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}

	adD := &io.SellerSaveDataMeta{}
	filterD := bson.M{"dataMetaHash": input.DataMetaHash, "ontId": ontId}
	err = sql.FindElt(sql.DataMetaCollection, filterD, adD)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}

	arr := strings.Split(ontId, ":")
	if len(arr) != 3 {
		return qrCode.QrCodeResponse{}, err
	}
	sellerAddress, err := common.AddressFromBase58(arr[2])

	// dataMeta related in data contract tx.
	tokenTemplate := param.TokenTemplate{
		DataIDs:   adD.DataIds,
		TokenHash: adT.TokenMetaHash,
	}
	//itemMetaHash, err := ddxf.HashObject(input.ItemMeta)
	itemMetaDataBs, err := json.Marshal(input.ItemMeta)
	bs := sha256.Sum256(itemMetaDataBs)
	itemMetaHash, err := common.Uint256ParseFromBytes(bs[:])

	im := sellerconfig.ItemMeta{
		ItemMetaHash: itemMetaHash.ToHexString(),
		ItemMetaData: input.ItemMeta,
	}
	err = sql.InsertElt(sql.ItemMetaCollection, im)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	resourceIdBytes, rosourceDDOBytes, itemBytes := contract.ConstructPublishParam(sellerAddress, tokenTemplate,
		adT.TokenEndpoint, itemMetaHash, adD.ResourceType, adD.Fee, adD.ExpiredDate, adD.Stock, adD.DataIds)
	qrCodex, err := qrCode2.BuildPublishQrCode(sellerconfig.DefSellerConfig.NetType, input.MPContractHash,
		resourceIdBytes, rosourceDDOBytes, itemBytes, arr[2], ontId)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}

	err = sql.InsertElt(sql.SellerQrCodeCollection, qrCodex)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	return qrCode2.BuildQrCodeResponse(qrCodex.QrCodeId), nil
}

func GetQrCodeByQrCodeIdService(qrCodeId string) (qrCode.QrCode, error) {
	filter := bson.M{"qrCodeId": qrCodeId}
	code := qrCode.QrCode{}
	err := sql.FindElt(sql.SellerQrCodeCollection, filter, &code)
	return code, err
}

func QrCodeCallBackService(param qrCode.QrCodeCallBackParam) error {
	filter := bson.M{"qrCodeId": param.ExtraData.Id}
	code := qrCode.QrCode{}
	err := sql.FindElt(sql.SellerQrCodeCollection, filter, &code)
	if err != nil {
		return err
	}
	uuidType := utils.UUIDType(code.QrCodeId)
	switch uuidType {
	case utils.UUID_TOKEN_SELLER_PUBLISH:
		resourceId, ddo, err := qrCode2.ParseFromBytes(code.QrCodeData)
		if err != nil {
			return err
		}
		adD := sellerconfig.ItemMeta{}
		filterD := bson.M{"dataMetaHash": ddo.DescHash}
		err = sql.FindElt(sql.ItemMetaCollection, filterD, &adD)
		if err != nil {
			return err
		}
		in := io.MPEndpointPublishItemMetaInput{
			SignedDDXFTx: param.SignedTx,
			ItemMeta: io.PublishItemMeta{
				OnchainItemID: resourceId,
				ItemMeta:      adD.ItemMetaData,
			},
		}
		output := DefSellerImpl.PublishMPItemMeta(in, param.ExtraData.OntId)
		return output.Error()
	default:
		return errors.New("wrong uuid type")
	}
}
