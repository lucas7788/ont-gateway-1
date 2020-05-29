package qrCode

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ontio/sagapi/sagaconfig"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"time"
)

type qrCodeDesc struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Price  string `json:"price"`
}

func BuildBuyGetQrCodeRsp(qrCodeId string) QrCodeResponse {
	get := GetQrCode{
		ONTAuthScanProtocol:"",
	}
	return QrCodeResponse{
		QrCode:get,
		Id:"",
	}
}

func BuildBuyQrCode(netType string, resourceId string, n int, buyer string) (QrCode, error) {
	exp := time.Now().Unix() + QrCodeExpire
	data := &QrCodeData{
		Action: "signTransaction",
		Params: QrCodeParam{
			InvokeConfig: InvokeConfig{
				ContractHash: "",
				Functions: []Function{
					Function{
						Operation: "buyDToken",
						Args: []Arg{
							Arg{
								Name:  "resource_id",
								Value: "string:" + resourceId,
							},
							Arg{
								Name:  "n",
								Value: n,
							},
							Arg{
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
		return QrCode{}, err
	}
	id := common.GenerateUUId()
	sig, err := sagaconfig.DefSagaConfig.OntIdAccount.Sign(databs)
	if err != nil {
		return QrCode{}, err
	}

	qrDesc := qrCodeDesc{
		Type:   "invoke ddxf contract",
		Detail: "buyDtoken",
	}

	qrDescIn, err := json.Marshal(qrDesc)
	if err != nil {
		return QrCode{}, err
	}

	return QrCode{
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
