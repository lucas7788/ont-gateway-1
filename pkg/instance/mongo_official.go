package instance

import (
	"fmt"
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	instanceMongoOfficial *mongo.Database
	lockMongoOfficial     sync.Mutex
)

// MongoOfficial is singleton for mongo.Database
func MongoOfficial() *mongo.Database {
	if instanceMongoOfficial != nil {
		return instanceMongoOfficial
	}

	lockMongoOfficial.Lock()
	defer lockMongoOfficial.Unlock()
	if instanceMongoOfficial != nil {
		return instanceMongoOfficial
	}

	config := config.Load().MongoConfig
	db, err := storage.NewMongoOfficial(&config)
	if err != nil {
		panic(fmt.Sprintf("official mongo instantiate err:%v", err))
	}
	instanceMongoOfficial = db

	return instanceMongoOfficial
}
