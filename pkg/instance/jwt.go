package instance

import (
	"fmt"
	"sync"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/misc"
	"sync/atomic"
	"unsafe"
)

var (
	jwtMu  sync.Mutex
	jwtPtr unsafe.Pointer
)

// JWT is singleton for misk.JWT
func JWT() *misc.JWT {
	inst := atomic.LoadPointer(&jwtPtr)
	if inst != nil {
		return (*misc.JWT)(inst)
	}
	jwtMu.Lock()
	defer jwtMu.Unlock()

	inst = atomic.LoadPointer(&jwtPtr)
	if inst != nil {
		return (*misc.JWT)(inst)
	}
	jwt, err := misc.NewJWT(time.Hour*24, "HS256", []byte("ont-gateway-sec"))
	if err != nil {
		panic(fmt.Sprintf("NewJWT err:%v", err))
	}
	atomic.StorePointer(&jwtPtr, unsafe.Pointer(jwt))
	return jwt
}
