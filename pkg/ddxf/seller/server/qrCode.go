package server

import (
	"encoding/hex"
	"encoding/json"
	"github.com/kataras/go-errors"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ontio/ontology/common"
	common2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"strings"
	"time"
)

func BuildQrCodeResponse(id string) qrCode.QrCodeResponse {
	return qrCode.QrCodeResponse{
		QrCode: qrCode.GetQrCode{
			ONTAuthScanProtocol: ONTAuthScanProtocol + "/" + id,
		},
		Id: id,
	}
}
func ParsePublishParamFromQrCodeData(qrCodeData string) (resourceId []byte, resourceDdo []byte, item []byte, err error) {
	data := qrCode.QrCodeData{}
	err = json.Unmarshal([]byte(qrCodeData), &data)
	if err != nil {
		return
	}
	args := data.Params.InvokeConfig.Functions[0].Args
	if len(args) != 3 {
		err = errors.New("")
		return
	}
	if args[0].Name != "resourceId" {
		err = errors.New("")
		return
	}
	resourceIdStr := strings.ReplaceAll(args[0].Value.(string), "ByteArray:", "")
	resourceId, err = hex.DecodeString(resourceIdStr)
	ddoStr := strings.ReplaceAll(args[1].Value.(string), "ByteArray:", "")
	itemStr := strings.ReplaceAll(args[2].Value.(string), "ByteArray:", "")
	resourceDdo, err = hex.DecodeString(ddoStr)
	item, err = hex.DecodeString(itemStr)
	return
}

func ParseFromBytes(qrCodeData string) (resourceId []byte, resourceDdo *market_place_contract.ResourceDDO, item *market_place_contract.DTokenItem, err error) {
	var ddoBytes, itemBytes []byte
	resourceId, ddoBytes, itemBytes, err = ParsePublishParamFromQrCodeData(qrCodeData)
	if err != nil {
		return
	}
	source := common.NewZeroCopySource(ddoBytes)
	resourceDdo = &market_place_contract.ResourceDDO{}
	err = resourceDdo.Deserialize(source)
	if err != nil {
		return
	}
	item = &market_place_contract.DTokenItem{}
	err = item.FromBytes(itemBytes)
	return
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
								Value: "ByteArray:" + hex.EncodeToString(resource_ddo_bytes),
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
	id := common2.GenerateUUId(config.UUID_PUBLISH_ID)
	sig, err := ServerAccount.Sign(databs)
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
		Requester:  ServerAccount.Address.ToHexString(),
		Signer:     ontid,
		QrCodeData: string(databs),
		Callback:   QrCodeCallback,
		Exp:        exp,
		Chain:      chain,
		QrCodeDesc: string(qrDescIn),
	}, nil
}
