package service

import (
	"fmt"
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
			fmt.Printf("**********txHash:%s, err:%s addr:%s \n", input.SignedTx, err, addr)
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			continue
		}
		mutTx, err := tx.IntoMutable()
		if err != nil {
			fmt.Printf("**********txHash:%s, err:%s addr:%s \n", input.SignedTx, err, addr)
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			continue
		}
		txHash, err := kit.SendTransaction(mutTx)
		if err != nil {
			fmt.Printf("**********txHash:%s, err:%s addr:%s \n", input.SignedTx, err, addr)
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			continue
		}
		output.TxHash = txHash.ToHexString()
	}

	return
}
