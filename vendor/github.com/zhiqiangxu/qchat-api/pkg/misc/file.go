package misc

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/zhiqiangxu/qrpc"
)

// requireFile will panic if err happens
func requireFile(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("ReadFile err:%v", err))
	}
	return qrpc.String(data)
}

// RequireFile combines multiple files into one
// joins files with new line
func RequireFile(files ...string) string {
	var content []string
	for _, file := range files {
		content = append(content, requireFile(file))
	}
	return strings.Join(content, "\r\n")
}
