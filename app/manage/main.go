package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/app/manage/cmd"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "init_gw",
				Usage:  "add mongo index and bootstrap data etc",
				Action: cmd.AddMongoIndex,
			},
			{
				Name:   "poll_notify_tx",
				Usage:  "poll and notify for transactions",
				Action: cmd.PollNotifyTx,
			},
			{
				Name:   "update_notify_payment",
				Usage:  "update payment balance and notify payment recharge",
				Action: cmd.UpdateNotifyPayment,
			},
			&cmd.WalletCmd,
			&cmd.AppCmd,
			{
				Name:   "print",
				Usage:  "print structs",
				Action: cmd.Print,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
