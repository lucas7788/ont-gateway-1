package service

import (
	"context"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"go.uber.org/zap"
)

const (
	batch = 10
)

// PollTx impl
func (gw *Gateway) PollTx(ctx context.Context) (output io.PollTxOutput) {
	kit := sdk.NewOntologySdk()

	for {
		select {
		case <-ctx.Done():
			output.Msg = ctx.Err().Error()
			return
		default:
		}

		txlist, err := model.TxManager().QueryToPoll(batch)
		if err != nil {
			logger.Instance().Error("QueryToPoll", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		if len(txlist) == 0 {
			logger.Instance().Info("PollTx txlist empty")
			time.Sleep(time.Second * 5)
			continue
		}

		for _, tx := range txlist {
			event, err := kit.GetSmartContractEvent(tx.Hash)
			if err != nil {
				logger.Instance().Error("GetSmartContractEvent", zap.Error(err))
				model.TxManager().UpdatePollError(tx.Hash, err.Error())
				continue
			}
			if event == nil {
				logger.Instance().Error("GetSmartContractEvent returns nil event")
				model.TxManager().UpdatePollError(tx.Hash, "nil event")
				continue
			}

			_, err = model.TxManager().UpdateResultAndState(tx.Hash, model.TxPollResultExists, model.TxStateToNotify)
			if err != nil {
				logger.Instance().Error("UpdateResultAndState", zap.Error(err))
			}
		}
	}

}
