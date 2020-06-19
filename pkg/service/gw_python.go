package service

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/util"
)

// Python impl
func (gw *Gateway) Python(input io.PythonInput) (output io.PythonOutput) {

	name, err := ioutil.TempDir("/tmp", "gw")
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	defer os.RemoveAll(name)

	err = ioutil.WriteFile(name, []byte(input.Py), 0644)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	cmd := exec.Command("python", name)

	out, err := cmd.Output()
	output.Out = util.String(out)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}

	return
}
