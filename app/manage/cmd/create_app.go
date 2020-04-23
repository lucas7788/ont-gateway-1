package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// CreateApp for create new app
func CreateApp(c *cli.Context) error {

	id, err := model.AppManager().GetMaxAppIDFromDB()

	apps, _ := model.AppManager().GetAllFromDB()
	fmt.Println("id", id, "err", err, "apps", apps)

	input := io.CreateAppInput{ID: id + 1, Name: c.String("name"), TxNotifyURL: c.String("txNotifyUrl"), PaymentNotifyURL: c.String("paymentNotifyUrl")}
	output := service.Instance().CreateApp(input)
	return output.Error()
}
