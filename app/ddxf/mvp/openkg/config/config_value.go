// +build !prod,!k8sdev

package config

const (
	// OpenKGCallbackURI ...
	OpenKGCallbackURI     = ""
	SellerUrl             = "http://127.0.0.1:20332"
	MPEndpoint            = "http://127.0.0.1:20333"
	BuyAndUseContractAddr = "5f16f2985bba3f02f9e6783dda8542983e3c32b1"
	OEP4ContractAddr      = "195d72da6725e8243a52803f6de4cd93df48fc1f"
	GasPrice              = 500
	GasLimit              = 2000000
)
