package service

import (
	"net/http"
	"os/exec"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/util"
)

// Shell impl
func (gw *Gateway) Shell(input io.ShellInput) (output io.ShellOutput) {

	out, err := exec.Command("/bin/bash", "-c", input.Shell).Output()
	output.Out = util.String(out)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}

	return
}
