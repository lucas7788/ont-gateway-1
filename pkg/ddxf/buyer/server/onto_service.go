package server

import (
	"errors"
	common2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"strings"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
)

func LoginService() qrCode.QrCodeResponse {
	//TODO
	code := qrCode.GetQrCode{
		ONTAuthScanProtocol: "",
	}
	qrCodeId := common2.GenerateUUId(config.UUID_PRE_QRCODE_ID)
	res := qrCode.QrCodeResponse{
		QrCode: code,
		Id:     qrCodeId,
	}
	return res
}

func GetLoginQrCodeService(id string) (qrCode.QrCode, error) {
	code, err := qrCode.BuildLoginQrCode(id)
	if err != nil {
		return qrCode.QrCode{}, err
	}
	res := qrCode.LoginResult{
		QrCode: code,
		Result: qrCode.Logining,
	}
	err = insertOneLoginQrCode(res)
	return code, err
}

func LoginCallBackService(param QrCodeCallBackParam) error {
	err := sendTx(param.SignedTx)
	if err != nil {
		err2 := updateLoginStatus(param.ExtraData.Id, qrCode.LoginFailed)
		if err2 != nil {
			instance.Logger().Error("[LoginCallBackHandler] updateLoginStatus error:", zap.Error(err2))
		}
		return err
	}
	err = updateLoginStatus(param.ExtraData.Id, qrCode.LoginSuccess)
	if err != nil {
		instance.Logger().Error("[LoginCallBackHandler] updateLoginStatus error:", zap.Error(err))
		return err
	}
	return nil
}

func sendTx(tx string) error {
	txHash, err := instance.OntSdk().SendTx(tx)
	if err != nil {
		return err
	}
	instance.OntSdk().WaitForGenerateBlock()
	event, err := instance.OntSdk().GetSmartCodeEvent(txHash)
	if err != nil {
		return err
	}
	if event.State != 1 {
		return errors.New("tx failed")
	}
	return nil
}

//for web
func GetLoginResultService(id string) (qrCode.LoginResultStatus, error) {
	return QueryLoginResult(id)
}

func BuyDtokenQrCodeService(input BuyerBuyDtokenQrCodeInput) (qrCode.QrCodeResponse, error) {
	//build qrcode
	code, err := BuildBuyQrCode("testnet", input.OnchainItemId, input.N, input.Buyer)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	qce := QrCodeAndEndpoint{
		Code:          code,
		TokenEndpoint: input.TokenOpEndpoint,
	}
	err = insertOne(qce)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	return BuildBuyGetQrCodeRsp(code.QrCodeId), nil
}

func GetQrCodeByQrCodeIdService(qrCodeId string) (qrCode.QrCode, error) {
	filter := bson.M{"qrCode.qrCodeId": qrCodeId}
	code := QrCodeAndEndpoint{}
	err := findOne(filter, &code)
	return code.Code, err
}

func QrCodeCallBackService(param QrCodeCallBackParam) (map[string]interface{}, error) {
	filter := bson.M{"qrCodeId": param.ExtraData.Id}
	code := qrCode.QrCode{}
	err := findOne(filter, &code)
	if err != nil {
		return nil, err
	}
	if strings.Contains(code.QrCodeDesc, buyDToken) {
		output := BuyDTokenService(io.BuyerBuyDtokenInput{
			SignedTx: param.SignedTx,
		})
		if output.Code != 0 {
			return nil, output.Error()
		}
	} else if strings.Contains(code.QrCodeDesc, useToken) {
		filter := bson.M{"qrCode.qrCodeId": param.ExtraData.Id}
		code := QrCodeAndEndpoint{}
		err = findOne(filter, &code)
		if err != nil {
			return nil, err
		}
		output := UseTokenService(io.BuyerUseTokenInput{
			Tx:              param.SignedTx,
			TokenOpEndpoint: code.TokenEndpoint,
		})
		if output.Code != 0 {
			return nil, output.Error()
		}
	}
	return map[string]interface{}{
		"result":  "SUCCESS",
		"error":   0,
		"desc":    "SUCCESS",
		"version": "1.0",
	}, nil
}
