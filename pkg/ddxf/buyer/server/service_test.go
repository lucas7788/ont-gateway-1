package server

import (
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
)

func TestBuyDTokenService(t *testing.T) {
	input := io.BuyerBuyDtokenInput{
		SignedTx: "00d2ce8bd45ef40100000000000000c2eb0b000000005c645c529cbab407589537ef4b87b84374f23c384bd4cdc2fffc31a9f281fed2174717548d777b26cf360962757944746f6b656e077265736f5f323501000000000000000000000000000000151dd3ecff3994999739bee170e6f490437248a7000241403711b300ec7650b95a29003675ab20f517b21d42d5d725b7babae4558304f00623163825589c27452df6e1c093d40aef8b51ff00854fc6241ef187ceaa40bbcc2321025af6199b152051fb7d508d11897f8e95fa4c95aa76f764dda347f59e9db82955ac4140a91bcdb00468687b3c0177c0a26a1ae2d3d8f055fe1a9cbfa42f6e228ff66ebf43bbfcacb7a4830627a71a3c49297d768c3e87fa33a16e96471c0bf886ae12472321034a0e9b2b5478145833be19b8ae687a8b4625930288ab56f32614537bd40b22a6ac",
	}
	output := BuyDTokenService(input)
	assert.Equal(t, 0, output.Code)
	fmt.Println("BuyDTokenService output: ", output)
}

func TestHandleEvent(t *testing.T) {
	res, err := HandleEvent("0f792177d846c2e4a69e0a7a2058ced610febf701e8a671a9b0cb4447a5e1416", "buyDtoken")
	assert.Nil(t,err)
	fmt.Println(res)
}