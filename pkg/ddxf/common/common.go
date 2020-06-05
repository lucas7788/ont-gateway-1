package common

import (
	"encoding/hex"
	"errors"
	"github.com/ontio/ontology/common"
	"github.com/satori/go.uuid"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/param"
)

func GenerateUUId(preFix string) string {
	return preFix + uuid.NewV4().String()
}

func HandleEvent(txHash string, method string) ([]io.EndpointToken, error) {
	event, err := instance.OntSdk().GetSmartCodeEvent(txHash)
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
		if notify.ContractAddress == misc.DDXF_CONTRACT_ADDRESS {
			states, ok := notify.States.([]interface{})
			if !ok || len(states) != 4 {
				return nil, errors.New("notify wrong")
			}
			if method == "buyDToken" {
				buyer, err = common.AddressFromBase58(states[3].(string))
				if err != nil {
					return nil, err
				}
				onchainItemId = states[1].(string)
				onchainItemIdByes, err := hex.DecodeString(onchainItemId)
				if err != nil {
					return nil, err
				}
				res, err := instance.OntSdk().DefaultDDXFContract().PreInvoke("getTokenTemplates", []interface{}{onchainItemIdByes})
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
				return tokenEndpoints,nil
			} else if method == "useToken" {
				buyer, err = common.AddressFromBase58(states[2].(string))
				if err != nil {
					return nil, err
				}
				onchainItemId = states[1].(string)
				tokenTemplateHex := states[3].(string)
				tokenTemplateBytes, err := hex.DecodeString(tokenTemplateHex)
				if err != nil {
					return nil, err
				}
				tt := param.TokenTemplate{}
				err = tt.FromBytes(tokenTemplateBytes)
				if err != nil {
					return nil, err
				}
				return []io.EndpointToken{
					io.EndpointToken{
						Token:io.Token{
							TokenTemplate:tt,
							Buyer:buyer,
							OnchainItemId:onchainItemId,
						},
					},
				},nil
			}
		}
	}
	return nil, errors.New("")
}
