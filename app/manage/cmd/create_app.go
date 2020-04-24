package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// CreateApp for create new app
func CreateApp(c *cli.Context) error {

	id, err := model.AppManager().GetMaxAppIDFromDB()
	if err != nil {
		return err
	}

	ak, sk := model.AppManager().GenerateAkSk()

	input := io.CreateAppInput{
		ID:               id + 1,
		Name:             c.String("name"),
		TxNotifyURL:      c.String("txNotifyUrl"),
		PaymentNotifyURL: c.String("paymentNotifyUrl"),
		Ak:               ak,
		Sk:               sk,
	}
	output := service.Instance().CreateApp(input)
	return output.Error()
}
