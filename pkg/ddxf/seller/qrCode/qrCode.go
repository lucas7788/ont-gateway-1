package qrCode

import (
	"encoding/hex"
	"encoding/json"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sellerconfig"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/utils"
	"time"
)

type QrCodeResponse struct {
	QrCodeScan QrCodeScan `json:"qrCode"`
	Id         string     `json:"id"`
}

type QrCodeScan struct {
	ONTAuthScanProtocol string `json:"ONTAuthScanProtocol"`
}

func BuildQrCodeResponse(id string) *QrCodeResponse {
	return &QrCodeResponse{
		QrCodeScan: QrCodeScan{
			ONTAuthScanProtocol: sellerconfig.DefSellerConfig.ONTAuthScanProtocol + "/" + id,
		},
		Id: id,
	}
}

func BuildPublishQrCode(chain string, contractHash string, resourceId []byte, resource_ddo_bytes []byte, item_bytes []byte, payer string, ontid string) (*qrCode.QrCode, error) {
	exp := time.Now().Unix() + 600
	data := &qrCode.QrCodeData{
		Action: "signTransaction",
		Params: qrCode.QrCodeParam{
			InvokeConfig: qrCode.InvokeConfig{
				ContractHash: contractHash,
				Functions: []qrCode.Function{
					qrCode.Function{
						Operation: "token_seller_publish",
						Args: []qrCode.Arg{
							qrCode.Arg{
								Name:  "resourceId",
								Value: "ByteArray:" + hex.EncodeToString(resourceId),
							},
							qrCode.Arg{
								Name:  "resource_ddo_bytes",
								Value: "Address:" + hex.EncodeToString(resource_ddo_bytes),
							},
							qrCode.Arg{
								Name:  "item_bytes",
								Value: hex.EncodeToString(item_bytes),
							},
						},
					},
				},
				Payer:    payer,
				GasLimit: 20000,
				GasPrice: 500,
			},
		},
	}
	databs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	id := utils.GenerateUUId(utils.UUID_TOKEN_SELLER_PUBLISH)
	sig, err := sellerconfig.DefSellerConfig.ServerAccount.Sign(databs)
	if err != nil {
		return nil, err
	}

	qrDesc := qrCode.QrCodeDesc{
		Type:   "invoke wasm type",
		Detail: "token_seller_publish",
	}

	qrDescIn, err := json.Marshal(qrDesc)
	if err != nil {
		return nil, err
	}

	return &qrCode.QrCode{
		Ver:        "1.0.0",
		QrCodeId:   id,
		Signature:  hex.EncodeToString(sig),
		Requester:  sellerconfig.DefSellerConfig.ServerAccount.Address.ToHexString(),
		Signer:     ontid,
		QrCodeData: string(databs),
		Callback:   sellerconfig.DefSellerConfig.QrCodeCallback,
		Exp:        exp,
		Chain:      chain,
		QrCodeDesc: string(qrDescIn),
	}, nil
}
