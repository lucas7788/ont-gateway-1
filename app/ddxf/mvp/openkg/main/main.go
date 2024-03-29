package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/core/types"
	"github.com/zhiqiangxu/ont-gateway/app/ddxf/mvp/openkg/server"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
)

// MVP for openkg
func main() {

	if false {
		bs, _ := hex.DecodeString("03e951ceb3ab9ad212584d2e246a2890167f2194a76e7dec615f3610fe891fc482")
		pubkey, _ := keypair.DeserializePublicKey(bs)
		addr := types.AddressFromPubKey(pubkey)
		fmt.Println(addr.ToBase58())
		return
	}
	common.ConsortiumAddr = "http://113.31.112.154:20336"
	server.StartServer()
}
