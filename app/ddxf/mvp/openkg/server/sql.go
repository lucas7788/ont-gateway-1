package server

import (
	"context"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	OpenKgPublishParamCollection = "open_kg_publish_collection"
	UserInfoCollection           = "user_info_collection"
)

// Init for this collection
func initDb() (err error) {
	err = initUniqueIndex("u-seller", "openkg_id", OpenKgPublishParamCollection)
	if err != nil {
		return
	}
	err = initUniqueIndex("u-user-id", "user_id", UserInfoCollection)
	return
}

func initUniqueIndex(name, key string, collectionName string) (err error) {
	opts := &options.IndexOptions{}
	opts.SetName(name)
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: key, Value: bsonx.Int32(1)}},
		Options: opts,
	}
	_, err = instance.MongoOfficial().Collection(collectionName).Indexes().CreateOne(context.Background(), index)
	return
}

func InsertElt(collectionName string, data interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := instance.MongoOfficial().Collection(collectionName).InsertOne(ctx, data)
	return err
}

func FindElt(collectionName string, filter bson.M, data interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return instance.MongoOfficial().Collection(collectionName).FindOne(ctx, filter).Decode(data)
}

func FindManyElt(collectionName string, filter bson.M, data interface{}) (err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cursor, err := instance.MongoOfficial().Collection(collectionName).Find(ctx, filter)
	if err != nil {
		return
	}

	err = cursor.All(context.Background(), data)
	return
}
