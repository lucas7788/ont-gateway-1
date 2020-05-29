package io

// GetWalletMinInput for input
type GetWalletMinInput struct {
	WalletName string `json:"wallet_name"`
}

// GetWalletMinOutput for output
type GetWalletMinOutput struct {
	BaseResp
	Content string `json:"content"`
	Exists  bool   `json:"exists"`
}
