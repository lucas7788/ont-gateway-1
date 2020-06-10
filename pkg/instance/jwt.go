package instance

import (
	"fmt"
	"sync"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
)

var (
	instanceJWT *misc.JWT
	jwtOnce     sync.Once
)

// JWT is singleton for misk.JWT
func JWT() *misc.JWT {
	jwtOnce.Do(func() {
		jwt, err := misc.NewJWT(time.Hour*24, "HS256", []byte("ont-gateway-sec"))
		if err != nil {
			panic(fmt.Sprintf("NewJWT err:%v", err))
		}
		instanceJWT = jwt
	})

	return instanceJWT
}
