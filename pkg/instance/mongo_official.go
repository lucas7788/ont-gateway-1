package instance

import (
	"fmt"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	mongoPtr unsafe.Pointer
	mongoMu  sync.Mutex
)

// MongoOfficial is singleton for mongo.Database
func MongoOfficial() *mongo.Database {
	inst := atomic.LoadPointer(&mongoPtr)
	if inst != nil {
		return (*mongo.Database)(inst)
	}
	mongoMu.Lock()
	defer mongoMu.Unlock()

	inst = atomic.LoadPointer(&mongoPtr)
	if inst != nil {
		return (*mongo.Database)(inst)
	}

	conf := config.Load().MongoConfig
	db, err := storage.NewMongoOfficial(&conf)
	if err != nil {
		panic(fmt.Sprintf("official mongo instantiate err:%v", err))
	}
	atomic.StorePointer(&mongoPtr, unsafe.Pointer(db))
	return db
}
