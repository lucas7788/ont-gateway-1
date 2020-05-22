package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// AppCmd for manage app
var AppCmd = cli.Command{
	Action: cli.ShowSubcommandHelp,
	Name:   "app",
	Usage:  "Manage app",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "create an app for tx polling etc",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
				&cli.StringFlag{Name: "txNotifyUrl", Aliases: []string{"tu"}},
				&cli.StringFlag{Name: "paymentNotifyUrl", Aliases: []string{"pu"}},
			},
			Action: createApp,
		},
		{
			Name:   "ls",
			Usage:  "list app",
			Action: listApp,
		},
	},
}

func createApp(c *cli.Context) error {

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

func listApp(c *cli.Context) error {
	output := service.Instance().ListApp()
	fmt.Println("apps", output.Apps)
	return output.Error()
}
