package server

import (
	"encoding/hex"
	"encoding/json"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"time"
)

type qrCodeDesc struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Price  string `json:"price"`
}

func BuildBuyGetQrCodeRsp(qrCodeId string) qrCode.QrCodeResponse {
	get := qrCode.GetQrCode{
		ONTAuthScanProtocol: "",
	}
	return qrCode.QrCodeResponse{
		QrCode: get,
		Id:     qrCodeId,
	}
}

func BuildBuyQrCode(netType string, resourceId string, n int, buyer string) (qrCode.QrCode, error) {
	exp := time.Now().Unix() + qrCode.QrCodeExpire
	data := &qrCode.QrCodeData{
		Action: "signTransaction",
		Params: qrCode.QrCodeParam{
			InvokeConfig: qrCode.InvokeConfig{
				ContractHash: "",
				Functions: []qrCode.Function{
					qrCode.Function{
						Operation: "buyDToken",
						Args: []qrCode.Arg{
							qrCode.Arg{
								Name:  "resource_id",
								Value: "string:" + resourceId,
							},
							qrCode.Arg{
								Name:  "n",
								Value: n,
							},
							qrCode.Arg{
								Name:  "buyer",
								Value: "address:" + buyer,
							},
						},
					},
				},
				Payer:    buyer,
				GasLimit: 20000,
				GasPrice: 500,
			},
		},
	}
	databs, err := json.Marshal(data)
	if err != nil {
		return qrCode.QrCode{}, err
	}
	id := common.GenerateUUId()
	sig, err := BuyerMgrAccount.Sign(databs)
	if err != nil {
		return qrCode.QrCode{}, err
	}

	qrDesc := qrCodeDesc{
		Type:   "invoke ddxf contract",
		Detail: "buyDtoken",
	}

	qrDescIn, err := json.Marshal(qrDesc)
	if err != nil {
		return qrCode.QrCode{}, err
	}

	return qrCode.QrCode{
		Ver:        "1.0.0",
		QrCodeId:   id,
		Requester:  buyer,
		Signature:  hex.EncodeToString(sig),
		Signer:     buyer,
		QrCodeData: string(databs),
		Callback:   "",
		Exp:        exp,
		Chain:      netType,
		QrCodeDesc: string(qrDescIn),
	}, nil
}
