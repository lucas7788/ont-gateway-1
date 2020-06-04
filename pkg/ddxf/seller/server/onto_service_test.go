package server

import (
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
)

func TestPublishMetaService(t *testing.T) {
	input := io.SellerPublishMPItemMetaInput{
		ItemMeta:       map[string]interface{}{},
		TokenMetaHash:  "",
		DataMetaHash:   "",
		MPContractHash: "",
		MPEndpoint:     "",
	}
	res, err := PublishMetaService(input, "ontid")
	assert.Nil(t, err)
	fmt.Println("res: ", res)
}
