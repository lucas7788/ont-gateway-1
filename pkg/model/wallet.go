package model

import (
	"encoding/base64"

	"github.com/zhiqiangxu/ont-gateway/pkg/crypto"
	"github.com/zhiqiangxu/util"
)

// Wallet model
type Wallet struct {
	Name          string `bson:"name" json:"name"`
	CipherContent string `bson:"cipher_content" json:"cipher_content"`
}

var (
	passWd = []byte("!ont-gateway-wallet!")
)

// SetPlainContent will convert plain text to cipher text
func (w *Wallet) SetPlainContent(content string) (err error) {
	aesBytes, err := crypto.AESEncrypt(util.Slice(content), passWd)
	if err != nil {
		return
	}

	w.CipherContent = base64.StdEncoding.EncodeToString(aesBytes)

	return
}

// GetPlainContent returns plain content
func (w *Wallet) GetPlainContent() (content string, err error) {
	base64Decoded, err := base64.StdEncoding.DecodeString(w.CipherContent)
	if err != nil {
		return
	}

	contentBytes, err := crypto.AESDecrypt(base64Decoded, passWd)
	if err != nil {
		return
	}

	content = util.String(contentBytes)
	return
}
