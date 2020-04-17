package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
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
			app, exists := model.AppManager().GetApp(tx.App)
			if !exists {
				logger.Instance().Error("NotifyTx App not exists", zap.String("txHash", tx.Hash), zap.Int("app", tx.App))
				model.TxManager().UpdateState(tx.Hash, model.TxStateDone)
				continue
			}

			err = gw.notifyTx(app.TxNotifyURL, tx.Hash, tx.Result)
			if err != nil {
				logger.Instance().Error("notifyTx", zap.Int("app", tx.App), zap.String("txHash", tx.Hash), zap.Error(err))
				model.TxManager().UpdateNotifyError(tx.Hash, err.Error())
				continue
			}

			model.TxManager().UpdateState(tx.Hash, model.TxStateDone)
		}

	}
}

type notifyTxInput struct {
	TxHash string
	Result model.TxPollResult
}

func (gw *Gateway) notifyTx(url, txHash string, result model.TxPollResult) (err error) {

	input := notifyTxInput{TxHash: txHash, Result: result}
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
