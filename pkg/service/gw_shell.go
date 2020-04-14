package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/util"
)

// Shell impl
func (gw *Gateway) Shell(input io.ShellInput) (output io.ShellOutput) {

	name, err := ioutil.TempDir("/tmp", "gw")
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	defer os.RemoveAll(name)

	cmd := exec.Command("/bin/bash", "-c", input.Shell)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("TMP_DIR=%s", name))

	out, err := cmd.Output()
	output.Out = util.String(out)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}

	return
}
