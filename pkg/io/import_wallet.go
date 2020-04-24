package io

import "fmt"

// ImportWalletInput for input
type ImportWalletInput struct {
	WalletName    string `json:"wallet_name"`
	Content       string `json:"content"`
	CipherContent string `json:"cipher_content"`
}

// ImportWalletOutput for output
type ImportWalletOutput struct {
	BaseResp
}

// Validate ImportWalletInput
func (input *ImportWalletInput) Validate() (err error) {
	switch {
	case input.WalletName == "":
		return fmt.Errorf("wallet_name empty")
	case input.Content == "" && input.CipherContent == "":
		return fmt.Errorf("both content/cipher_content empty")
	}
	return nil
}
