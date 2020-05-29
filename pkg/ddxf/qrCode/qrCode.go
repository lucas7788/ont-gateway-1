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
	QrCodeId     string `json:"id" bson:"qrCodeId"`
	Ver          string `json:"ver" bson:"ver"`
	OrderId      string `json:"orderId" bson:"orderId"`
	Requester    string `json:"requester" bson:"requester"`
	Signature    string `json:"signature" bson:"signature"`
	Signer       string `json:"signer" bson:"signer"`
	QrCodeData   string `json:"data" bson:"qrCodeData"`
	Callback     string `json:"callback" bson:"callback"`
	Exp          int64  `json:"exp" bson:"exp"`
	Chain        string `json:"chain" bson:"chain"`
	QrCodeDesc   string `json:"desc" bson:"qrCodeDesc"`
	ContractType string `json:"contractType" bson:"contractType"`
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
