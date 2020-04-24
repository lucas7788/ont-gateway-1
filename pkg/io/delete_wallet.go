package io

// DeleteWalletInput for input
type DeleteWalletInput struct {
	WalletName string `json:"wallet_name"`
}

// DeleteWalletOutput for output
type DeleteWalletOutput struct {
	BaseResp
}
