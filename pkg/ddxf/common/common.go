package common

import (
	"encoding/hex"
	"errors"
	"github.com/ontio/ontology/common"
	"github.com/satori/go.uuid"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
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
				break
			} else if method == "useToken" {
				buyer, err = common.AddressFromBase58(states[2].(string))
				if err != nil {
					return nil, err
				}
				onchainItemId = states[1].(string)
				break
			}
		}
	}

	onchainItemIdByes, err := hex.DecodeString(onchainItemId)
	if err != nil {
		return nil, err
	}
	res, err := instance.OntSdk().DDXFContract(0, 0,
		nil).PreInvoke("getTokenTemplates", []interface{}{onchainItemIdByes})
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
