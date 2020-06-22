package service

import (
	"net/http"

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
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			return
		}
		output.TxHash = txHash.ToHexString()
	}

	return
}
