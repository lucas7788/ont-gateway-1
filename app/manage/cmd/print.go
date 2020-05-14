package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// Print ...
func Print(c *cli.Context) (err error) {
	bytes, err := json.Marshal(io.GetResourceOutput{})
	if err != nil {
		return
	}

	fmt.Println(string(bytes))
	return
}
