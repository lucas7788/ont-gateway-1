package common

import (
	"encoding/hex"
	"errors"

	"github.com/ont-bizsuite/ddxf-sdk/ddxf_contract"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	uuid "github.com/satori/go.uuid"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
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
	for _, notify := range event.Notify {
		//TODO ddxf contractaddress
		if notify.ContractAddress == misc.DDXF_CONTRACT_ADDRESS {
			states, ok := notify.States.([]interface{})
			if !ok {
				return nil, errors.New("notify wrong")
			}
			if method == "buyDToken" {
				return handleBuyDtokenEvt(states)
			} else if method == "useToken" {
				return handleUseTokenEvt(states)
			} else if method == "buyAndUseToken" {
				handleBuyDtokenEvt(states)
				handleUseTokenEvt(states)
			}
		}
	}
	return nil, errors.New("")
}

func HanleBuyAndUseToken(txHash string) ([]io.EndpointToken, []io.EndpointToken, error) {
	event, err := instance.OntSdk().GetSmartCodeEvent(txHash)
	if err != nil {
		return nil, nil, err
	}
	if event.State != 1 {
		return nil, nil, errors.New("tx failed")
	}
	var buyEndpointToken []io.EndpointToken
	var useEndpointToken []io.EndpointToken
	for _, notify := range event.Notify {
		//TODO ddxf contractaddress
		if notify.ContractAddress == misc.DDXF_CONTRACT_ADDRESS {
			states, ok := notify.States.([]interface{})
			if !ok {
				return nil, nil, errors.New("notify wrong")
			}
			if states[0] == "buyDtoken" {
				buyEndpointToken, err = handleBuyDtokenEvt(states)
				if err != nil {
					return nil, nil, err
				}
			} else if states[0] == "" {
				useEndpointToken, err = handleUseTokenEvt(states)
				if err != nil {
					return nil, nil, err
				}
			}
		}
	}
	return buyEndpointToken, useEndpointToken, nil
}

func handleUseTokenEvt(states []interface{}) ([]io.EndpointToken, error) {
	var buyer common.Address
	var onchainItemId string
	var err error
	if len(states) != 5 {
		return nil, errors.New("event failed")
	}
	buyer, err = common.AddressFromBase58(states[2].(string))
	if err != nil {
		return nil, err
	}
	onchainItemId = states[1].(string)
	tokenTemplateHex := states[4].(string)
	tokenTemplateBytes, err := hex.DecodeString(tokenTemplateHex)
	if err != nil {
		return nil, err
	}
	tt := ddxf_contract.TokenTemplate{}
	err = tt.FromBytes(tokenTemplateBytes)
	if err != nil {
		return nil, err
	}
	return []io.EndpointToken{
		io.EndpointToken{
			Token: io.Token{
				TokenTemplate: tt,
				Buyer:         buyer,
				OnchainItemId: onchainItemId,
			},
		},
	}, nil
}

func handleBuyDtokenEvt(states []interface{}) ([]io.EndpointToken, error) {
	var buyer common.Address
	var onchainItemId string
	var err error
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
	return tokenEndpoints, nil
}

var ConsortiumAddr string

// SendTx for send to main-net and Consortium
func SendTx(tx string) (txHash string, err error) {
	addrs := []string{misc.GetOntNode()}
	if ConsortiumAddr != "" {
		addrs = append(addrs, ConsortiumAddr)
	}
	input := io2.SendTxInput{SignedTx: tx, Addrs: addrs}
	output := service.Instance().SendTx(input)
	txHash, err = output.TxHash, output.Error()
	return
}

// SendRawTx for send to main-net and Consortium
func SendRawTx(tx *types.MutableTransaction) (txHash string, err error) {
	addrs := []string{misc.GetOntNode()}
	if ConsortiumAddr != "" {
		addrs = append(addrs, ConsortiumAddr)
	}
	input := io2.SendRawTxInput{Tx: tx, Addrs: addrs}
	output := service.Instance().SendRawTx(input)
	txHash, err = output.TxHash, output.Error()
	return
}
