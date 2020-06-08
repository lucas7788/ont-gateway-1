package seller_buyer

import (
	"encoding/json"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"fmt"
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
