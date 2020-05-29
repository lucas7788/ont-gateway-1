package server

import (
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"errors"
)

func BuyDtokenQrCodeService(input io.BuyerBuyDtokenInput) (qrCode.QrCodeResponse, error) {
	//build qrcode
	code, err := qrCode.BuildBuyQrCode("testnet", input.OnchainItemId, input.N, input.Buyer)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	err = insertOne(code)
	if err != nil {
		return qrCode.QrCodeResponse{}, err
	}
	return qrCode.BuildBuyGetQrCodeRsp(code.QrCodeId), nil
}

func GetQrCodeByQrCodeIdService(qrCodeId string) (qrCode.QrCode, error) {
	filter := bson.M{"qrCodeId": qrCodeId}
	code := qrCode.QrCode{}
	err := findOne(filter, &code)
	return code, err
}

func QrCodeCallBackService(param QrCodeCallBackParam) (map[string]interface{}, error) {
	filter := bson.M{"qrCodeId": param.ExtraData.Id}
	code := qrCode.QrCode{}
	err := findOne(filter, &code)
	if err != nil {
		return nil, err
	}
	var method string
	if strings.Contains(code.QrCodeDesc, buyDToken) {
		method = buyDToken
	} else if strings.Contains(code.QrCodeDesc, useToken) {
		method = useTokenM
	}
	endpointTokens, err := sendTxAndGetTokens(param.SignedTx, method)
	err = insertOne(endpointTokens)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"result":  "SUCCESS",
		"error":   0,
		"desc":    "SUCCESS",
		"version": "1.0",
	}, nil
}

func sendTxAndGetTokens(txHex string, method string) ([]io.EndpointToken, error) {
	tx, err := utils.TransactionFromHexString(txHex)
	if err != nil {
		return nil, err
	}
	mutTx, err := tx.IntoMutable()
	if err != nil {
		return nil, err
	}
	txHash, err := instance.OntSdk().GetKit().SendTransaction(mutTx)
	if err != nil {
		return nil, err
	}
	event, err := instance.OntSdk().GetSmartCodeEvent(txHash.ToHexString())
	if err != nil {
		return nil, err
	}
	if event.State != 1 {
		return nil, errors.New("tx failed")
	}
	var buyer common.Address
	var onchainItemId string
	for _, notify := range event.Notify {
		//TODO ddxf contractaddress
		if notify.ContractAddress == "" {
			states, ok := notify.States.([]string)
			if !ok || len(states) != 4 {
				return nil, errors.New("notify wrong")
			}
			if method == "buyDtoken" {
				buyer, err = common.AddressFromBase58(states[3])
				if err != nil {
					return nil, err
				}
				onchainItemId = states[1]
			} else if method == "useToken" {
				buyer, err = common.AddressFromBase58(states[2])
				if err != nil {
					return nil, err
				}
				onchainItemId = states[1]
			}
		}
	}

	res, err := instance.OntSdk().DDXFContract(0, 0,
		nil).PreInvoke("getTokenTemplates", []interface{}{onchainItemId})
	if err != nil {
		return nil, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return nil, err
	}
	tokenEndpoints, err := io.ConstructTokensAndEndpoint(data, buyer, onchainItemId)
	if err != nil {
		return nil, err
	}
	return tokenEndpoints, nil
}
