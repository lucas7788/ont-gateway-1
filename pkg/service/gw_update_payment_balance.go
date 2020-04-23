package service

import (
	"context"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	paymentBatch = 50
)

// UpdatePaymentBalance impl
func (gw *Gateway) UpdatePaymentBalance(ctx context.Context) (output io.UpdatePaymentBalanceOutput) {

	client := instance.MongoOfficial().Client()

	for {
		select {
		case <-ctx.Done():
			output.Msg = ctx.Err().Error()
			return
		default:
		}

		payments, err := model.PaymentManager().QueryBalanceExpired(paymentBatch)
		if err != nil {
			logger.Instance().Error("QueryBalanceExpired", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}
		if len(payments) == 0 {
			logger.Instance().Info("QueryBalanceExpired payments empty")
			time.Sleep(time.Second * 5)
			continue
		}

		for i := range payments {
			payment := payments[i]
			err = client.UseSession(ctx, func(sessionContext mongo.SessionContext) (err error) {
				err = sessionContext.StartTransaction()
				if err != nil {
					return
				}
				defer sessionContext.AbortTransaction(sessionContext)

				err = gw.updateBalance(sessionContext, &payment, nil)
				if err != nil {
					return
				}

				err = sessionContext.CommitTransaction(sessionContext)
				return
			})
			if err != nil {
				logger.Instance().Error("UseSession", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}
		}
	}
}
