package main

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

// MVP for openkg
func main() {
	r := gin.Default()
	r.POST(generateOntIdByUserIdURI, generateOntIdByUserId)
	r.POST(publishURI, publish)
	r.POST(buyAndUseURI, buyAndUse)

	wallet, err := instance.OntSdk().GetKit().OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Println("OpenWallet failed: ", err)
		return
	}
	payer, err = wallet.GetAccountByAddress("", []byte("123456"))
	if err != nil {
		fmt.Println("GetAccountByAddress failed: ", err)
		return
	}
	instance.DDXFSdk().SetPayer(payer)
	defPlainSeed = hex.EncodeToString(sha256.New().Sum(keypair.SerializePrivateKey(payer.PrivateKey)))
	r.Run(":" + openkgPort)
}
