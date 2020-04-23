package cmd

import (
	"context"

	"github.com/oklog/run"
	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
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
		logger.Instance().Error("UpdatePaymentBalance output", zap.Any("output", output))
		return nil
	}, func(err error) {
		logger.Instance().Error("UpdatePaymentBalance interrupt", zap.Error(err))
		cancelFunc()
	})

	g.Add(func() error {
		output := service.Instance().NotifyPaymentBalance(ctx)
		logger.Instance().Error("NotifyPaymentBalance output", zap.Any("output", output))
		return nil
	}, func(err error) {
		logger.Instance().Error("NotifyPaymentBalance interrupt", zap.Error(err))
		cancelFunc()
	})

	return g.Run()
}
