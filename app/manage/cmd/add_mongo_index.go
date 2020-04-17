package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// AddMongoIndex for add mongo index
func AddMongoIndex(c *cli.Context) error {

	err := model.AddonConfigManager().EnsureIndex()
	fmt.Println("AddonConfigManager.EnsureIndex", err)
	err = model.AddonDeploymentManager().EnsureIndex()
	fmt.Println("AddonDeploymentManager.EnsureIndex", err)

	err = model.TxManager().EnsureIndex()
	fmt.Println("TxManager.EnsureIndex", err)

	return nil
}
