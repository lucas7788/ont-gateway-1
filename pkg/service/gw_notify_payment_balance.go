package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"github.com/zhiqiangxu/util"
	"go.uber.org/zap"
)

// NotifyPaymentBalance impl
func (gw *Gateway) NotifyPaymentBalance(ctx context.Context) (output io.NotifyPaymentBalanceOutput) {

	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			output.Msg = ctx.Err().Error()
			return
		default:
		}

		nothing2do := true

		payments, err := model.PaymentManager().QueryToNotifyPreRecharging(paymentBatch)
		if err != nil {
			instance.Logger().Error("QueryToNotifyPreRecharging", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		if len(payments) > 0 {
			nothing2do = false
		}

		for i := range payments {
			payment := &payments[i]
			util.GoFunc(&wg, func() {
				err := gw.notifyPaymentBalance(payment)
				if err != nil {
					_, err = model.PaymentManager().UpdateNotifyError(payment.App, payment.PaymentID, err.Error())
					if err != nil {
						instance.Logger().Error("UpdateNotifyError", zap.Error(err))
					}
				}
			})
		}
		wg.Wait()

		payments, err = model.PaymentManager().QueryToNotifyRecharging(paymentBatch)
		if err != nil {
			instance.Logger().Error("QueryToNotifyRecharging", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		if len(payments) > 0 {
			nothing2do = false
		}

		for i := range payments {
			payment := &payments[i]
			util.GoFunc(&wg, func() {
				err := gw.notifyPaymentBalance(payment)
				if err != nil {
					_, err = model.PaymentManager().UpdateNotifyError(payment.App, payment.PaymentID, err.Error())
					if err != nil {
						instance.Logger().Error("UpdateNotifyError", zap.Error(err))
					}
				}
			})
		}
		wg.Wait()

		if nothing2do {
			instance.Logger().Info("NotifyPaymentBalance nothing to do")
			time.Sleep(time.Second * 5)
			continue
		}

	}

}

type notifyPaymentBalanceInput struct {
	App       int             `json:"app"`
	PaymentID string          `json:"payment_id"`
	Balance   int             `json:"balance"`
	StartTime time.Time       `json:"start_time"`
	PayPeriod model.PayPeriod `json:"pay_period"`
}

func (gw *Gateway) notifyPaymentBalance(payment *model.Payment) (err error) {
	app := model.AppManager().GetByID(payment.App)
	if app == nil {
		instance.Logger().Error("notifyPaymentBalance App not exists", zap.String("paymentID", payment.PaymentID), zap.Int("app", payment.App))
		err = fmt.Errorf("App not exists")
		return
	}

	input := notifyPaymentBalanceInput{App: payment.App, PaymentID: payment.PaymentID, Balance: payment.Balance, StartTime: payment.StartTime, PayPeriod: payment.PayPeriod}
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return
	}

	_, _, body, err := forward.PostJSONRequest(app.PaymentNotifyURL, jsonBytes, nil)
	if err != nil {
		return
	}

	if !bytes.Equal(body, []byte("ok")) {
		err = fmt.Errorf("invalid notify resp:%s", util.String(body))
		return
	}

	return
}
