package io

import "fmt"

// ImportWalletInput for input
type ImportWalletInput struct {
	WalletName string `json:"wallet_name"`
	Content    string `json:"content"`
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
	case input.Content == "":
		return fmt.Errorf("content empty")
	}
	return nil
}
