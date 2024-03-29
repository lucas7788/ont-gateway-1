package service

import (
	"fmt"
	"net/http"
	"os/exec"
	"path"
	"time"

	"github.com/kballard/go-shellquote"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// DumpWallet impl
func (gw *Gateway) DumpWallet(input io.DumpWalletInput) (output io.DumpWalletOutput) {

	conf := config.Load()
	out := time.Now().Format("2006_01_02")
	fmt.Println("mongodump", "--collection=wallet --out="+path.Join(input.Path, out)+" --uri="+shellquote.Join(conf.MongoConfig.ConnectionString))
	cmd := exec.Command("mongodump", "--collection", "wallet", "--out", path.Join(input.Path, out), "--uri", conf.MongoConfig.ConnectionString)

	_, err := cmd.Output()
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	return
}
