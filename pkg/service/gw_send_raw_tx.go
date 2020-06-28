package service

import (
	"net/http"

	"fmt"

	osdk "github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// SendRawTx impl
func (gw *Gateway) SendRawTx(input io.SendRawTxInput) (output io.SendRawTxOutput) {

	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	for _, addr := range input.Addrs {
		kit := osdk.NewOntologySdk()
		kit.NewRpcClient().SetAddress(addr)
		txHash, err := kit.SendTransaction(input.Tx)

		if err != nil {
			fmt.Printf("**********txHash:%s, err:%s addr:%s \n", txHash.ToHexString(), err, addr)
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			continue
		}
		output.TxHash = txHash.ToHexString()
	}

	return
}
