package io

// GetWalletInput for input
type GetWalletInput struct {
	WalletName string `json:"wallet_name"`
}

// GetWalletOutput for output
type GetWalletOutput struct {
	BaseResp
	Content string `json:"content"`
	Exists  bool   `json:"exists"`
}
