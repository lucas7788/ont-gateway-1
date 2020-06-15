package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// Print ...
func Print(c *cli.Context) (err error) {
	ak, _ := model.AppManager().GenerateAkSk()
	fmt.Println(ak)
	return

	bytes, err := json.Marshal(io.GetResourceOutput{})
	if err != nil {
		return
	}

	fmt.Println(string(bytes))
	return
}
