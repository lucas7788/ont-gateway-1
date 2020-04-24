package cmd

import (
	"context"

	"github.com/oklog/run"
	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
	"go.uber.org/zap"
)

// UpdateNotifyPayment for update and notify payment
func UpdateNotifyPayment(c *cli.Context) error {
	var g run.Group

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	g.Add(func() error {
		output := service.Instance().UpdatePaymentBalance(ctx)
		instance.Logger().Error("UpdatePaymentBalance output", zap.Any("output", output))
		return nil
	}, func(err error) {
		instance.Logger().Error("UpdatePaymentBalance interrupt", zap.Error(err))
		cancelFunc()
	})

	g.Add(func() error {
		output := service.Instance().NotifyPaymentBalance(ctx)
		instance.Logger().Error("NotifyPaymentBalance output", zap.Any("output", output))
		return nil
	}, func(err error) {
		instance.Logger().Error("NotifyPaymentBalance interrupt", zap.Error(err))
		cancelFunc()
	})

	return g.Run()
}
