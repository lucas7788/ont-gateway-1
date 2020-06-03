package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// AddMongoIndex for add mongo index
func AddMongoIndex(c *cli.Context) error {

	err := model.AddonConfigManager().Init()
	fmt.Println("AddonConfigManager.Init", err)
	err = model.AddonDeploymentManager().Init()
	fmt.Println("AddonDeploymentManager.Init", err)

	err = model.TxManager().Init()
	fmt.Println("TxManager.Init", err)
	err = model.AppManager().Init()
	fmt.Println("AppManager.Init", err)

	err = model.PaymentConfigManager().Init()
	fmt.Println("PaymentConfigManager.Init", err)
	err = model.PaymentManager().Init()
	fmt.Println("PaymentManager.Init", err)
	err = model.PaymentOrderManager().Init()
	fmt.Println("PaymentOrderManager.Init", err)
	err = model.ResourceVersionManager().Init()
	fmt.Println("ResourceVersionManager.Init", err)

	err = model.WalletManager().Init()
	fmt.Println("WalletManager.Init", err)

	return nil
}
