package storage

import (
	"context"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// NewMongoOfficial for official mongo client
func NewMongoOfficial(conf *config.MongoConfig) (db *mongo.Database, err error) {

	var client *mongo.Client
	{
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		client, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.ConnectionString))
		if err != nil {
			return
		}
	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			return
		}
	}

	cs, err := connstring.Parse(conf.ConnectionString)
	if err != nil {
		return
	}

	db = client.Database(cs.Database)

	return
}
