package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
)

var test bool

// CreatePaymentOrder impl
func (gw *Gateway) CreatePaymentOrder(input io.CreatePaymentOrderInput) (output io.CreatePaymentOrderOutput) {
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

		paymentConfig, err := model.PaymentConfigManager().GetWithTx(sessionContext, input.App, input.PaymentConfigID)
		if err != nil {
			return
		}

		if paymentConfig == nil {
			err = fmt.Errorf("paymentConfig %v not exist", input.PaymentConfigID)
			return
		}

		if !paymentConfig.HasCoinType(input.CoinType) {
			err = fmt.Errorf("coin type %v not supported", input.CoinType)
			return
		}
		if !paymentConfig.HasPayMethod(input.PayMethod) {
			err = fmt.Errorf("pay method %v not supported", input.PayMethod)
			return
		}

		paidAmount, err := model.PaymentOrderManager().TotalAmountPaidWithTx(sessionContext, input.App, input.PaymentID)
		if err != nil {
			return
		}

		payment, err := model.PaymentManager().GetOnePayment(sessionContext, input.App, input.PaymentID)
		if err != nil {
			return
		}

		var toBillAmount int
		if payment == nil {

			unitAmount, exists := paymentConfig.AmountForPeriod(input.PayPeriod)
			if !exists {
				err = fmt.Errorf("non existing PayPeriod %v for PaymentConfig %v ||", input.PayPeriod, paymentConfig.PaymentConfigID)
				return
			}

			payment := model.Payment{App: input.App, PaymentID: input.PaymentID, PaymentConfigID: input.PaymentConfigID, PayPeriod: input.PayPeriod, PayMethod: input.PayMethod, UnitAmount: unitAmount}
			toBillAmount, err = payment.AmountForNth(0, paymentConfig)
			if err != nil {
				return
			}

			payment.Balance = input.Amount - toBillAmount
			payment.UpdatedAt = time.Now()
			if payment.Balance >= 0 {
				payment.State = model.PaymentStateStarted
				payment.StartTime = payment.UpdatedAt
				payment.BalanceExpireTime = payment.StartTime.Add(payment.PeriodDuration())
			} else {
				payment.State = model.PaymentStateToStart
			}
			output.Balance = payment.Balance

			err = model.PaymentManager().Insert(sessionContext, payment)
			if err != nil {
				return
			}

		} else {

			if payment.IsEnded() {
				err = fmt.Errorf("payment %v is ended at %v", input.PaymentID, payment.EndTime)
				return
			}

			var nth int
			switch payment.IsStarted() {
			case true:
				nth = payment.PeriodNth()
			case false:
				nth = 0
			}

			toBillAmount, err = payment.AmountForNth(nth, paymentConfig)
			if err != nil {
				return
			}

			output.Balance = input.Amount + paidAmount - toBillAmount

			switch payment.IsStarted() {
			case true:
				balanceExpireTime := payment.StartTime.Add(time.Duration(nth+1) * (payment.PeriodDuration()))
				_, err = model.PaymentManager().UpdateBalanceAndExpireTime(sessionContext, input.App, input.PaymentID, output.Balance, balanceExpireTime)
			case false:
				if output.Balance >= 0 {
					startTime := time.Now()
					balanceExpireTime := startTime.Add(time.Duration(nth+1) * (payment.PeriodDuration()))
					_, err = model.PaymentManager().UpdateBalanceAndStartTime(sessionContext, input.App, input.PaymentID, output.Balance, startTime, balanceExpireTime)
				} else {
					_, err = model.PaymentManager().UpdateBalance(sessionContext, input.App, input.PaymentID, output.Balance)
				}
			}

			if err != nil {
				return
			}

		}

		order := model.PaymentOrder{App: input.App, PaymentID: input.PaymentID, OrderID: input.OrderID, Amount: input.Amount, CoinType: input.CoinType, OrderInfo: input.OrderInfo}
		err = model.PaymentOrderManager().InsertWithTx(sessionContext, order)
		if err != nil {
			return
		}

		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			return
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
