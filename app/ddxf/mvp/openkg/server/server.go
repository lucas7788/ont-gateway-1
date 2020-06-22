package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
)

const (
	openkgPort               = "10999"
	publishURI               = "/publish"
	buyAndUseURI             = "/buyAndUse"
	generateOntIdByUserIdURI = "/generateOntIdByUserId"
)

var (
	payer        *ontology_go_sdk.Account
	defPlainSeed string
	defGasPrice  = uint64(500)
	defGasLimit  = uint64(20000000)
)

func InitData() error {
	wallet, err := instance.OntSdk().GetKit().OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Println("OpenWallet failed: ", err)
		return err
	}
	payer, err = wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", []byte("123456"))
	if err != nil {
		fmt.Println("GetAccountByAddress failed: ", err)
		return err
	}
	instance.DDXFSdk().SetPayer(payer)
	defPlainSeed = hex.EncodeToString(sha256.New().Sum(keypair.SerializePrivateKey(payer.PrivateKey)))
	return nil
}

// MVP for openkg
func StartServer() {
	if err := InitData(); err != nil {
		fmt.Println("init data failed: ", err)
		return
	}
	r := gin.Default()
	r.POST(generateOntIdByUserIdURI, GenerateOntIdByUserId)
	r.POST(publishURI, Publish)
	r.POST(buyAndUseURI, BuyAndUse)
	r.Run(":" + openkgPort)
}
