package service

import (
	"net/http"

	osdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// SendTx impl
func (gw *Gateway) SendTx(input io.SendTxInput) (output io.SendTxOutput) {

	err := input.Validate()
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	for _, addr := range input.Addrs {
		kit := osdk.NewOntologySdk()
		kit.NewRpcClient().SetAddress(addr)
		tx, err := utils.TransactionFromHexString(input.SignedTx)
		if err != nil {
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			return
		}
		mutTx, err := tx.IntoMutable()
		if err != nil {
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			return
		}
		txHash, err := kit.SendTransaction(mutTx)
		if err != nil {
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			return
		}
		output.TxHash = txHash.ToHexString()
	}

	return
}
