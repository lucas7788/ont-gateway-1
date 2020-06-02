package misc

import (
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"gotest.tools/assert"
)

func TestChain(t *testing.T) {

	sdk := NewOntSdk()
	amount, err := sdk.GetAmountTransferred("ff250fdf93f4e888b1f1f94af9792fad84d64406fe2ae409204f844b6876b1f1")
	assert.Assert(t, err == nil && amount == 10000000)

	name := "/tmp/wfile"
	wallet, err := sdk.GetKit().CreateWallet(name)
	assert.Assert(t, err == nil, err)

	defer os.RemoveAll(name)

	u1 := uuid.NewV4()
	psw := u1.String()
	_, err = wallet.NewDefaultSettingAccount([]byte(psw))
	assert.Assert(t, err == nil)
	err = wallet.Save()
	assert.Assert(t, err == nil)

}
