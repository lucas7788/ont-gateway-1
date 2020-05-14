package model

import (
	"context"
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// ResourceVersionMgr ...
type ResourceVersionMgr struct {
}

var (
	resourceVersionMgr     *ResourceVersionMgr
	resourceVersionMgrLock sync.Mutex
)

// ResourceVersionManager ...
func ResourceVersionManager() *ResourceVersionMgr {
	if resourceVersionMgr != nil {
		return resourceVersionMgr
	}

	resourceVersionMgrLock.Lock()
	defer resourceVersionMgrLock.Unlock()

	if resourceVersionMgr != nil {
		return resourceVersionMgr
	}

	resourceVersionMgr = &ResourceVersionMgr{}
	return resourceVersionMgr
}

const (
	resourceVersionCollectionName = "resource_version"
)

// UpdateResource ...
func (m *ResourceVersionMgr) UpdateResource(rv ResourceVersion) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = instance.MongoOfficial().Collection(resourceVersionCollectionName).InsertOne(ctx, rv)
	if err != nil {
		return
	}

	if we, ok := err.(mongo.WriteException); ok {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				exists = true
				err = nil
			}
		}
	}
	return
}

// ForceUpdateResource ...
func (m *ResourceVersionMgr) ForceUpdateResource(rv ResourceVersion) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": rv.App, "id": rv.ID, "block": rv.Block}

	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "desc", Value: rv.Desc},
			bson.E{Key: "hash", Value: rv.Hash}},
		}}

	opt := options.Update().SetUpsert(true)

	result, err := instance.MongoOfficial().Collection(resourceVersionCollectionName).UpdateOne(ctx, filter, update, opt)
	if err != nil {
		return
	}

	exists = result.MatchedCount > 0
	return
}

// DeleteResourceByID is only called for test purpose
func (m *ResourceVersionMgr) DeleteResourceByID(app int, id string) (n int64, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": app, "id": id}

	result, err := instance.MongoOfficial().Collection(resourceVersionCollectionName).DeleteMany(ctx, filter)
	if err != nil {
		return
	}

	n = result.DeletedCount
	return
}

// GetByBlock ...
func (m *ResourceVersionMgr) GetByBlock(app int, id string, block uint32) (rv *ResourceVersion, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": app, "id": id, "block": block}

	rv = &ResourceVersion{}
	err = instance.MongoOfficial().Collection(resourceVersionCollectionName).FindOne(ctx, filter).Decode(rv)
	if err == mongo.ErrNoDocuments {
		rv = nil
		err = nil
		return
	}
	if err != nil {
		return
	}

	return
}

// GetByHash ...
func (m *ResourceVersionMgr) GetByHash(app int, id, hash string) (rv *ResourceVersion, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": app, "id": id, "hash": hash}

	rv = &ResourceVersion{}
	err = instance.MongoOfficial().Collection(resourceVersionCollectionName).FindOne(ctx, filter).Decode(rv)
	if err == mongo.ErrNoDocuments {
		rv = nil
		err = nil
		return
	}
	if err != nil {
		return
	}

	return
}

// Init for this collection
func (m *ResourceVersionMgr) Init() (err error) {

	blockOpts := &options.IndexOptions{}
	blockOpts.SetName("u-app-id-block")
	blockOpts.SetUnique(true)
	blockIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "app", Value: bsonx.Int32(1)}, {Key: "id", Value: bsonx.Int32(1)}, {Key: "block", Value: bsonx.Int32(1)}},
		Options: blockOpts,
	}

	hashOpts := &options.IndexOptions{}
	hashOpts.SetName("u-app-id-hash")
	hashOpts.SetUnique(true)
	hashIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "app", Value: bsonx.Int32(1)}, {Key: "id", Value: bsonx.Int32(1)}, {Key: "hash", Value: bsonx.Int32(1)}},
		Options: hashOpts,
	}

	models := []mongo.IndexModel{
		blockIndex,
		hashIndex,
	}
	_, err = instance.MongoOfficial().Collection(resourceVersionCollectionName).Indexes().CreateMany(context.Background(), models)
	return
}
