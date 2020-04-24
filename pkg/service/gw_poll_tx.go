package service

import (
	"context"
	"sync"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"github.com/zhiqiangxu/util"
	"go.uber.org/zap"
)

const (
	batch = 50
)

// PollTx impl
func (gw *Gateway) PollTx(ctx context.Context) (output io.PollTxOutput) {

	kit := instance.OntSdk().GetKit()

	var wg sync.WaitGroup

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

		for i := range txlist {
			tx := txlist[i]

			util.GoFunc(&wg, func() {
				event, err := kit.GetSmartContractEvent(tx.Hash)
				if err != nil {
					logger.Instance().Error("GetSmartContractEvent", zap.Error(err))

					if tx.IsExpired() {
						model.TxManager().FinishPoll(tx.Hash, model.TxPollResultExpired, err.Error())
					} else {
						model.TxManager().UpdatePollError(tx.Hash, err.Error())
					}
				}
				if event == nil {
					logger.Instance().Error("GetSmartContractEvent returns nil event")

					errMsg := "nil event"
					if tx.IsExpired() {
						model.TxManager().FinishPoll(tx.Hash, model.TxPollResultExpired, errMsg)
					} else {
						model.TxManager().UpdatePollError(tx.Hash, errMsg)
					}
				}

				_, err = model.TxManager().FinishPoll(tx.Hash, model.TxPollResultExists, "")
				if err != nil {
					logger.Instance().Error("FinishPoll", zap.Error(err))
				}
			})

		}

		wg.Wait()
	}

}
