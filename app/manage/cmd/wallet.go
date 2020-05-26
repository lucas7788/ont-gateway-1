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
			ArgsUsage: "wallet_name wallet_psw wallet_address",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "ciphered", Aliases: []string{"c"}},
			},
		},
		{
			Name:      "get",
			Usage:     "get a wallet",
			Action:    getWallet,
			ArgsUsage: "wallet_name",
		},
		{
			Name:      "delete",
			Usage:     "delete a wallet",
			Action:    deleteWallet,
			ArgsUsage: "wallet_name",
		},
	},
}

func importWallet(c *cli.Context) (err error) {

	data, err := ioutil.ReadFile(c.Args().Get(2))
	if err != nil {
		return
	}

	input := io.ImportWalletInput{PSW: c.Args().Get(1), WalletName: c.Args().Get(0)}
	if false && c.Bool("ciphered") {
		// input.CipherContent = string(data)
	} else {
		input.Content = string(data)
	}

	output := service.Instance().ImportWallet(input)
	err = output.Error()
	return
}

func getWallet(c *cli.Context) (err error) {
	input := io.GetWalletInput{WalletName: c.Args().Get(0)}
	output := service.Instance().GetWallet(input)
	err = output.Error()
	return
}

func deleteWallet(c *cli.Context) (err error) {

	input := io.DeleteWalletInput{WalletName: c.Args().Get(0)}
	output := service.Instance().DeleteWallet(input)
	err = output.Error()
	return
}
