package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
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

func importWallet(c *cli.Context) error {

	fmt.Println(c.Args().Get(0), c.Args().Get(1))
	return nil
}
