package model

import (
	"encoding/base64"

	"github.com/zhiqiangxu/ont-gateway/pkg/crypto"
	"github.com/zhiqiangxu/util"
)

// Wallet model
type Wallet struct {
	Name          string `bson:"name" json:"name"`
	CipherPSW     string `bson:"cipher_psw" json:"cipher_psw"`
	CipherContent string `bson:"cipher_content" json:"cipher_content"`
}

var (
	passWd = crypto.SHA256([]byte("!!ont-gw-wallet!"))
)

// SetPlain will convert plain text to cipher text
func (w *Wallet) SetPlain(psw, content string) (err error) {

	aesBytes, err := crypto.AESEncrypt(util.Slice(content), passWd)
	if err != nil {
		return
	}

	w.CipherContent = base64.StdEncoding.EncodeToString(aesBytes)

	aesBytes, err = crypto.AESEncrypt(util.Slice(psw), passWd)
	if err != nil {
		return
	}

	w.CipherPSW = base64.StdEncoding.EncodeToString(aesBytes)
	return
}

// GetPlain returns plain content
func (w *Wallet) GetPlain() (psw, content string, err error) {
	base64Decoded, err := base64.StdEncoding.DecodeString(w.CipherContent)
	if err != nil {
		return
	}

	contentBytes, err := crypto.AESDecrypt(base64Decoded, passWd)
	if err != nil {
		return
	}

	content = util.String(contentBytes)

	base64Decoded, err = base64.StdEncoding.DecodeString(w.CipherPSW)
	if err != nil {
		return
	}

	pswBytes, err := crypto.AESDecrypt(base64Decoded, passWd)
	if err != nil {
		return
	}

	psw = util.String(pswBytes)
	return
}
