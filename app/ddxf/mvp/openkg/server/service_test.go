package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	if err := InitData(); err != nil {
		fmt.Println("err: ", err)
		return
	}
	m.Run()
}

func TestGenerateOntIdService(t *testing.T) {
	input := GenerateOntIdInput{
		ReqID:  "req_id",
		UserId: "use_id",
	}
	output := GenerateOntIdService(input)
	fmt.Println(output.Msg)
	fmt.Println("output:", output)
	assert.Equal(t, output.Code, 0)
}

func TestPublishService(t *testing.T) {
	input2 := PublishInput{
		ReqID:    "req_id",
		OpenKGID: "",
		UserID:   "use_id",
		Item: map[string]interface{}{
			"item": "val",
		},
		Datas:     []map[string]interface{}{},
		OnChainId: "",
	}
	output2 := PublishService(input2)
	fmt.Println("output2:", output2)
	assert.Equal(t, output2.Code, 0)
}
