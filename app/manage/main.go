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
				Name:   "add_mongo_index",
				Usage:  "add mongo index",
				Action: cmd.AddMongoIndex,
			},
			{
				Name:   "poll_notify_tx",
				Usage:  "poll and notify for transactions",
				Action: cmd.PollNotifyTx,
			},
			{
				Name:  "create_app",
				Usage: "create an app for tx polling etc",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
					&cli.StringFlag{Name: "txNotifyUrl", Aliases: []string{"u"}},
				},
				Action: cmd.CreateApp,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
