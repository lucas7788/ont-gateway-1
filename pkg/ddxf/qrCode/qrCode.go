package qrCode

import (
	"encoding/hex"
	"encoding/json"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"time"
)

const (
	QrCodeExpire = 10 * 60
	testnet      = "testnet"
)

type QrCodeCallBackParam struct {
	Signer    string    `json:"signer"`
	SignedTx  string    `json:"signedTx"`
	ExtraData ExtraData `json:"extraData"`
}

type ExtraData struct {
	Id        string `json:"id"`
	PublicKey string `json:"publickey"`
	OntId     string `json:"ontId"`
}

type QrCodeDesc struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Price  string `json:"price"`
}

func BuildLoginQrCode(id string) (QrCode, error) {
	qp := QrCodeParam{
		InvokeConfig: InvokeConfig{
			ContractHash: "",
			Functions: []Function{
				Function{
					Operation: "signMessage",
					Args: []Arg{Arg{
						Name:  "login",
						Value: "String:ontid",
					}},
				},
			},
			Payer:    "",
			GasPrice: 500,
			GasLimit: 20000,
		},
	}
	qcd := QrCodeData{
		Action: "signMessage",
		Params: qp,
	}
	qcdBs, err := json.Marshal(qcd)
	if err != nil {
		return QrCode{}, err
	}
	sig, err := config.DefDDXFConfig().OperatorAccount.Sign(qcdBs)
	if err != nil {
		return QrCode{}, err
	}
	qrCode := QrCode{
		QrCodeId:     id,
		Ver:          "v2.0.0",
		Requester:    config.DefDDXFConfig().OperatorOntid,
		Signature:    hex.EncodeToString(sig),
		Signer:       "",
		QrCodeData:   string(qcdBs),
		Callback:     "", //TODO
		Exp:          time.Now().Unix() + QrCodeExpire,
		Chain:        testnet,
		QrCodeDesc:   "",
		ContractType: "",
	}
	return qrCode, nil
}

type LoginResultStatus uint8

const (
	NotLogin LoginResultStatus = iota
	Logining
	LoginFailed
	LoginSuccess
)

type LoginResult struct {
	QrCode QrCode            `json:"qrCode" bson:"qrCode"`
	Result LoginResultStatus `json:"result" bson:"result"`
}

type QrCode struct {
	QrCodeId     string `json:"id" bson:"qrCodeId"`
	Ver          string `json:"ver" bson:"ver"`
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

type QrCodeResponse struct {
	QrCode GetQrCode `json:"qrCode"`
	Id     string    `json:"id"`
}

type GetQrCode struct {
	ONTAuthScanProtocol string `json:"ONTAuthScanProtocol"`
}
