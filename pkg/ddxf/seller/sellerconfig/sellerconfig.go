package sellerconfig

import (
	"encoding/json"
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/common/password"
	"io/ioutil"
	"os"
)

const (
	DEFAULT_LOG_LEVEL = log.InfoLog
	DEFAULT_REST_PORT = uint(8080)
)

const (
	ONT_MAIN_NET = "http://dappnode1.ont.io:20336"
	ONT_TEST_NET = "http://polaris1.ont.io:20336"
	ONT_SOLO_NET = "http://127.0.0.1:20336"
)

const (
	NETWORK_ID_SOLO_NET int = 1
	NETWORK_ID_TEST_NET int = 2
	NETWORK_ID_MAIN_NET int = 3
)

const (
	TestNetType = "TestNet"
	MainNetType = "MainNet"
	SoloNetType = "soloNet"
)

type SellerConfig struct {
	NetWorkId           int    `json:"network_id"`
	Version             string `json:"version"`
	RestPort            uint   `json:"rest_port"`
	ONTAuthScanProtocol string `json:"ontauth_scan_protocol"`
	QrCodeCallback      string `json:"qrcode_callback"`
	NetType             string `json:"net_type"`
	WalletName          string `json:"wallet_name"`
	OntSdk              *sdk.OntologySdk
	ServerAccount       *sdk.Account
}

var DefSellerConfig = &SellerConfig{
	NetWorkId:           NETWORK_ID_SOLO_NET,
	Version:             "1.0",
	RestPort:            DEFAULT_REST_PORT,
	ONTAuthScanProtocol: "http://172.29.36.101/ddxf/seller/getQrCodeDataByQrCodeId",
	QrCodeCallback:      "http://172.29.36.101/ddxf/seller/qrCodeSendTx",
	NetType:             SoloNetType,
}

func InitAccount(s *sdk.OntologySdk, WalletName string) (*Account, error) {
	wallet, err := s.OpenWallet(WalletName)
	if err != nil {
		return nil, fmt.Errorf("error in OpenWallet:%s\n", err)
	}

	passwd, err := password.GetAccountPassword()
	if err != nil {
		return nil, fmt.Errorf("input password error %s", err)
	}

	ServerAccount, err := wallet.GetDefaultAccount(passwd)
	if err != nil {
		return nil, fmt.Errorf("error in GetDefaultAccount:%s\n", err)
	}

	return ServerAccount, nil
}

func InitSellerConfig(configFileName string) (*SellerConfig, error) {
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return DefSellerConfig, nil
	}
	file, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	cfg := &SellerConfig{}
	err = json.Unmarshal(bs, cfg)
	if err != nil {
		return nil, err
	}
	s := sdk.NewOntologySdk()
	switch cfg.NetWorkId {
	case NETWORK_ID_SOLO_NET:
		s.NewRpcClient().SetAddress(ONT_SOLO_NET)
		cfg.NetType = SoloNetType
	case NETWORK_ID_MAIN_NET:
		s.NewRpcClient().SetAddress(ONT_MAIN_NET)
		cfg.NetType = MainNetType
	case NETWORK_ID_TEST_NET:
		s.NewRpcClient().SetAddress(ONT_TEST_NET)
		cfg.NetType = TestNetType
	default:
		return nil, fmt.Errorf("Errof NetType %d", cfg.NetWorkId)
	}

	account, err := InitAccount(s, cfg.WalletName)
	if err != nil {
		return nil, err
	}
	cfg.ServerAccount = account

	return cfg, nil
}
