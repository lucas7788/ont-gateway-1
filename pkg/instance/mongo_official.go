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
	mongoOnce sync.Once
)

// MongoOfficial is singleton for mongo.Database
func MongoOfficial() *mongo.Database {
	mongoOnce.Do(func() {
		conf := config.Load().MongoConfig
		db, err := storage.NewMongoOfficial(&conf)
		if err != nil {
			panic(fmt.Sprintf("official mongo instantiate err:%v", err))
		}
		instanceMongoOfficial = db
	})

	return instanceMongoOfficial
}
