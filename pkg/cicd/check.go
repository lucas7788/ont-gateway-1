package cicd

import (
	"bytes"

	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
)

// Check url for readiness
func Check(url string) (ok bool, err error) {

	_, _, body, err := forward.Get(url)
	if err != nil {
		return
	}

	ok = bytes.Compare(body, []byte("ok")) == 0
	return
}
