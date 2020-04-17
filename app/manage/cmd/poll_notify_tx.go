package cmd

import (
	"context"

	"github.com/oklog/run"
	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
	"go.uber.org/zap"
)

// PollNotifyTx for polling and notifying tx
func PollNotifyTx(c *cli.Context) error {
	var g run.Group

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	g.Add(func() error {
		output := service.Instance().PollTx(ctx)
		logger.Instance().Error("PollTx output", zap.Any("output", output))
		return nil
	}, func(err error) {
		logger.Instance().Error("PollTx interrupt", zap.Error(err))
		cancelFunc()
	})

	g.Add(func() error {
		output := service.Instance().NotifyTx(ctx)
		logger.Instance().Error("NotifyTx output", zap.Any("output", output))
		return nil
	}, func(err error) {
		logger.Instance().Error("NotifyTx interrupt", zap.Error(err))
		cancelFunc()
	})

	return g.Run()
}
