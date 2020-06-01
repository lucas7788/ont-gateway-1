package service

import (
	"errors"
	utils2 "github.com/ontio/ontology-go-sdk/utils"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sql"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/utils"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
)

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
		tx, err := utils2.TransactionFromHexString(param.SignedTx)
		if err != nil {
			return err
		}
		mutTx, err := tx.IntoMutable()
		if err != nil {
			return err
		}
		txHash, err := instance.OntSdk().GetKit().SendTransaction(mutTx)
		if err != nil {
			return err
		}
		event, err := instance.OntSdk().GetSmartCodeEvent(txHash.ToHexString())
		if err != nil {
			return err
		}
		if event.State != 1 {
			return errors.New("tx failed")
		}

		return nil
	default:
		return errors.New("wrong uuid type")
	}
}
