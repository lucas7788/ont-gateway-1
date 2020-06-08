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
		TokenMetaHash:  "e2a740fa12bd94f0e242688e29f6d803f7671eb1f81bcfbdc1c3e213878e7dd4",
		DataMetaHash:   "f1dfe4c60f9f8e4942559ee14c549ce63abfced6a1be08519761744f2429ac35",
		MPContractHash: "",
		MPEndpoint:     "",
	}
	res, err := PublishMetaService(input, "did:ont:AcVBV1zKGogf9Q54p1Ve78NSQVU5ZUUGkn")
	assert.Nil(t, err)
	fmt.Println("res: ", res)
}
