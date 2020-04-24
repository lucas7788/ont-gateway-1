package cmd

import (
	"io/ioutil"

	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/service"
)

// WalletCmd for manage wallet
var WalletCmd = cli.Command{
	Action: cli.ShowSubcommandHelp,
	Name:   "wallet",
	Usage:  "Manage wallet",
	Subcommands: []*cli.Command{
		{
			Name:      "import",
			Usage:     "import a new wallet",
			Action:    importWallet,
			ArgsUsage: "wallet_name wallet_address",
		},
	},
}

func importWallet(c *cli.Context) (err error) {

	data, err := ioutil.ReadFile(c.Args().Get(1))
	if err != nil {
		return
	}

	output := service.Instance().ImportWallet(io.ImportWalletInput{WalletName: c.Args().Get(0), Content: string(data)})
	err = output.Error()
	return
}
