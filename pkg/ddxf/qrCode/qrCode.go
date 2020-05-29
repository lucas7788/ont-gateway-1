package qrCode

import ()

type SendTxParam struct {
	Signer    string    `json:"signer"`
	SignedTx  string    `json:"signedTx"`
	ExtraData ExtraData `json:"extraData"`
}

type ExtraData struct {
	Id        string `json:"id"`
	PublicKey string `json:"publickey"`
	OntId     string `json:"ontId"`
}

type QrCode struct {
	QrCodeId     string `json:"id" db:"QrCodeId"`
	Ver          string `json:"ver" db:"Ver"`
	OrderId      string `json:"orderId" db:"OrderId"`
	Requester    string `json:"requester" db:"Requester"`
	Signature    string `json:"signature" db:"Signature"`
	Signer       string `json:"signer" db:"Signer"`
	QrCodeData   string `json:"data" db:"QrCodeData"`
	Callback     string `json:"callback" db:"Callback"`
	Exp          int64  `json:"exp" db:"Exp"`
	Chain        string `json:"chain" db:"Chain"`
	QrCodeDesc   string `json:"desc" db:"QrCodeDesc"`
	ContractType string `json:"contractType" db:"ContractType"`
}

type QrCodeData struct {
	Action string      `json:"action"`
	Params QrCodeParam `json:"params"`
}

type QrCodeParam struct {
	InvokeConfig InvokeConfig `json:"invokeConfig"`
}

type InvokeConfig struct {
	ContractHash string     `json:"contractHash"`
	Functions    []Function `json:"functions"`
	Payer        string     `json:"payer"`
	GasLimit     uint64     `json:"gasLimit"`
	GasPrice     uint64     `json:"gasPrice"`
}

type Function struct {
	Operation string `json:"operation"`
	Args      []Arg  `json:"args"`
}

type Arg struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func BuildWetherForcastQrCode() (*QrCode, error) {
	return nil, nil
}
