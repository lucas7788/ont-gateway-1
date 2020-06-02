package service

import (
	"io/ioutil"
	"net/http"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
)

// CreateWallet impl
func (gw *Gateway) CreateWallet(input io.CreateWalletInput) (output io.CreateWalletOutput) {

	if input.WalletName == "" {
		input.WalletName = uuid.NewV4().String()
	}

	path := "/tmp/" + uuid.NewV4().String()
	wallet, err := instance.OntSdk().GetKit().CreateWallet(path)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	defer os.RemoveAll(path)

	u1 := uuid.NewV4()
	psw := u1.String()
	_, err = wallet.NewDefaultSettingAccount([]byte(psw))
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	err = wallet.Save()
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	w := model.Wallet{Name: input.WalletName}

	err = w.SetPlain(psw, string(data))
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	err = model.WalletManager().Insert(w)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	return
}
