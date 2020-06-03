package server

import (
	"context"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/qrCode"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	buyerCollectionName = "buyer"
	loginCollectionName = "buyer_login_qr_code"
)

func initDb() error {
	opts := &options.IndexOptions{}
	opts.SetName("u-qrCodeId")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "qrCodeId", Value: bsonx.Int32(1)}},
		Options: opts,
	}
	_, err := instance.MongoOfficial().Collection(buyerCollectionName).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return err
	}
	_, err = instance.MongoOfficial().Collection(loginCollectionName).Indexes().CreateOne(context.Background(), index)
	return err
}

func insertOneLoginQrCode(data interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := instance.MongoOfficial().Collection(loginCollectionName).InsertOne(ctx, data)
	return err
}

func updateLoginStatus(id string, res qrCode.LoginResultStatus) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	filter := bson.D{{"qrCodeId", id}}
	update := bson.D{{"result", res}}
	_, err := instance.MongoOfficial().Collection(loginCollectionName).UpdateOne(ctx, filter, update)
	return err
}

func QueryLoginResult(id string) (qrCode.LoginResultStatus, error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	filter := bson.M{"qrCode.qrCodeId": id}
	res := qrCode.LoginResult{}
	err := instance.MongoOfficial().Collection(loginCollectionName).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return qrCode.NotLogin, err
	}
	return res.Result, nil
}

func insertOne(data interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := instance.MongoOfficial().Collection(buyerCollectionName).InsertOne(ctx, data)
	return err
}

func insertMany(data []interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := instance.MongoOfficial().Collection(buyerCollectionName).InsertMany(ctx, data)
	return err
}

func findOne(filter bson.M, data interface{}) error {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return instance.MongoOfficial().Collection(buyerCollectionName).FindOne(ctx, filter).Decode(data)
}
