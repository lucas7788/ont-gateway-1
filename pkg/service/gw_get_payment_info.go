package service

import (
	"context"
	"net/http"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetPaymentInfo impl
func (gw *Gateway) GetPaymentInfo(input io.GetPaymentInfoInput) (output io.GetPaymentInfoOutput) {
	if !test {
		err := input.Validate()
		if err != nil {
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			return
		}
	}

	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := instance.MongoOfficial().Client()
	err := client.UseSession(ctx, func(sessionContext mongo.SessionContext) (err error) {
		err = sessionContext.StartTransaction()
		if err != nil {
			return
		}
		defer sessionContext.AbortTransaction(sessionContext)

		payment, err := model.PaymentManager().GetOnePayment(sessionContext, input.App, input.PaymentID)
		if err != nil {
			return
		}
		if payment == nil {
			return
		}

		output.Payment = payment

		paymentConfig, err := model.PaymentConfigManager().GetWithTx(sessionContext, input.App, payment.PaymentConfigID)
		if err != nil {
			return
		}
		if paymentConfig == nil {
			return
		}

		if payment.IsStarted() && payment.BalanceExpireTime.Before(time.Now()) {
			err = gw.updateBalance(sessionContext, payment, paymentConfig)
			if err != nil {
				return
			}
		}

		output.PaymentConfig = paymentConfig

		if input.WithOrders {
			output.PaymentOrders, err = model.PaymentOrderManager().GetPaymentOrders(sessionContext, input.App, input.PaymentID)
			if err != nil {
				return
			}
		}

		return
	})

	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	return
}

func (gw *Gateway) updateBalance(sessionContext mongo.SessionContext, payment *model.Payment, paymentConfig *model.PaymentConfig) (err error) {
	paidAmount, err := model.PaymentOrderManager().TotalAmountPaidWithTx(sessionContext, payment.App, payment.PaymentID)
	if err != nil {
		return
	}

	nth := payment.PeriodNth()
	toBillAmount, err := payment.AmountForNth(nth, paymentConfig)
	if err != nil {
		return
	}

	balance := paidAmount - toBillAmount
	balanceExpireTime := payment.StartTime.Add(time.Duration(nth+1) * (payment.PeriodDuration()))
	_, err = model.PaymentManager().UpdateBalanceAndExpireTime(sessionContext, payment.App, payment.PaymentID, balance, balanceExpireTime)
	return
}
