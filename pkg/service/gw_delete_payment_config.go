package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeletePaymentConfig impl
func (gw *Gateway) DeletePaymentConfig(input io.DeletePaymentConfigInput) (output io.DeletePaymentConfigOutput) {

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

		n, err := model.PaymentManager().CountForPaymentConfigWithTx(sessionContext, input.App, input.PaymentConfigID)
		if err != nil {
			return
		}
		if n > 0 {
			err = fmt.Errorf("can't delete payment config with existing payments")
			return
		}

		exists, err := model.PaymentConfigManager().DeleteOneWithTx(sessionContext, input.App, input.PaymentConfigID)
		if err != nil {
			return
		}

		err = sessionContext.CommitTransaction(sessionContext)

		output.Exists = exists

		return
	})
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}

	return
}
