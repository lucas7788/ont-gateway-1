package instance

import (
	"fmt"
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	instanceMongoOfficialClient *mongo.Database
	lockMongoOfficialClient     sync.Mutex
)

// MongoOfficialClient is singleton for mongo.Database
func MongoOfficialClient() *mongo.Database {
	if instanceMongoOfficialClient != nil {
		return instanceMongoOfficialClient
	}

	lockMongoOfficialClient.Lock()
	defer lockMongoOfficialClient.Unlock()
	if instanceMongoOfficialClient != nil {
		return instanceMongoOfficialClient
	}

	config := config.Load().MongoConfig
	db, err := storage.NewMongoOfficial(&config)
	if err != nil {
		panic(fmt.Sprintf("official mongo instantiate err:%v", err))
	}
	instanceMongoOfficialClient = db

	return instanceMongoOfficialClient
}
