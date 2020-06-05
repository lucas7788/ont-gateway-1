package storage

import (
	"context"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	StorageSaveDataMap = "StorageSaveDataMapCollection"
)

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
