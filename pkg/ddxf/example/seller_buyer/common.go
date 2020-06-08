package seller_buyer

import (
	"encoding/json"
	"fmt"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
)

func SendPOST(url string, param interface{}) ([]byte, error) {
	bs, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	_, _, data, err := forward.PostJSONRequest(url, bs)
	if err != nil {
		return nil, err
	}
	fmt.Println("data:", string(data))
	return data, nil
}
