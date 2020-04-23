package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"github.com/zhiqiangxu/util"
	"go.uber.org/zap"
)

// NotifyTx impl
func (gw *Gateway) NotifyTx(ctx context.Context) (output io.NotifyTxOutput) {
	for {
		select {
		case <-ctx.Done():
			output.Msg = ctx.Err().Error()
			return
		default:
		}

		txlist, err := model.TxManager().QueryToNotify(batch)
		if err != nil {
			logger.Instance().Error("QueryToNotify", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		if len(txlist) == 0 {
			logger.Instance().Info("NotifyTx txlist empty")
			time.Sleep(time.Second * 5)
			continue
		}
		for _, tx := range txlist {
			if tx.App == 0 {
				err = gw.notifyAdminTx(tx.Hash, tx.Result, tx.PollAmount)
				if err != nil {
					logger.Instance().Error("notifyAdminTx", zap.Int("app", tx.App), zap.String("txHash", tx.Hash), zap.Error(err))
					model.TxManager().UpdateNotifyError(tx.Hash, err.Error())
					continue
				}
			} else {
				app, exists := model.AppManager().GetApp(tx.App)
				if !exists {
					logger.Instance().Error("NotifyTx App not exists", zap.String("txHash", tx.Hash), zap.Int("app", tx.App))
					model.TxManager().UpdateState(tx.Hash, model.TxStateDone)
					continue
				}

				err = gw.notifyTx(app.TxNotifyURL, tx.Hash, tx.Result, tx.PollAmount)
				if err != nil {
					logger.Instance().Error("notifyTx", zap.Int("app", tx.App), zap.String("txHash", tx.Hash), zap.Error(err))
					model.TxManager().UpdateNotifyError(tx.Hash, err.Error())
					continue
				}
			}

			model.TxManager().UpdateState(tx.Hash, model.TxStateDone)
		}

	}
}

type notifyTxInput struct {
	TxHash   string             `json:"tx_hash"`
	NodeAddr string             `json:"node_addr"`
	Result   model.TxPollResult `json:"result"`
	Amount   uint64             `json:"amount"`
}

func (gw *Gateway) notifyAdminTx(hash string, result model.TxPollResult, pollAmount bool) (err error) {

	return
}

func (gw *Gateway) notifyTx(url, txHash string, result model.TxPollResult, pollAmount bool) (err error) {

	input := notifyTxInput{TxHash: txHash, NodeAddr: instance.OntSdkInstance().GetOntNode(), Result: result}
	if result == model.TxPollResultExists && pollAmount {
		var amount uint64
		amount, err = instance.OntSdkInstance().GetAmountTransferred(txHash)
		if err != nil {
			return
		}
		input.Amount = amount
	}
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return
	}
	_, _, body, err := forward.PostJSONRequest(url, jsonBytes)
	if err != nil {
		return
	}

	if !bytes.Equal(body, []byte("ok")) {
		err = fmt.Errorf("invalid notify resp:%s", util.String(body))
		return
	}

	return
}
