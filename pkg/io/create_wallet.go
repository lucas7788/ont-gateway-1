package io

// CreateWalletInput for input
type CreateWalletInput struct {
	WalletName string `json:"wallet_name"`
}

// CreateWalletOutput for output
type CreateWalletOutput struct {
	BaseResp
	WalletName string `json:"wallet_name"`
	Content    string `json:"content"`
	PSW        string `json:"psw"`
}
