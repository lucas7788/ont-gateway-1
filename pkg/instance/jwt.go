package instance

import (
	"fmt"
	"sync"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
)

var (
	instanceJWT *misc.JWT
	lockJWT     sync.Mutex
)

// JWT is singleton for misk.JWT
func JWT() *misc.JWT {
	if instanceJWT != nil {
		return instanceJWT
	}

	lockJWT.Lock()
	defer lockJWT.Unlock()
	if instanceJWT != nil {
		return instanceJWT
	}

	instanceJWT, err := misc.NewJWT(time.Hour*24, "HS256", []byte("ont-gateway-sec"))
	if err != nil {
		panic(fmt.Sprintf("NewJWT err:%v", err))
	}

	return instanceJWT
}
