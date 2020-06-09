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
	SellerQrCodeCollection = "seller_qr_code_collection"
	DataMetaCollection     = "seller_data_meta_collection"
	TokenMetaCollection    = "seller_token_meta_collection"
	ItemMetaCollection     = "seller_item_meta_collection"
	PublishParamCollection = "publish_param_collection"
)

// Init for this collection
func initDb() (err error) {
	opts := &options.IndexOptions{}
	opts.SetName("u-seller")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "seller", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	_, err = instance.MongoOfficial().Collection(sellerCollectionName).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return
	}
	index = mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "dataMetaHash", Value: bsonx.Int32(1)}},
		Options: opts,
	}
	_, err = instance.MongoOfficial().Collection(DataMetaCollection).Indexes().CreateOne(context.Background(), index)
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
