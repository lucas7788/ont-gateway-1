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

// PaymentConfigMgr for PaymentConfig crud
type PaymentConfigMgr struct {
}

var (
	paymentConfigMgr     *PaymentConfigMgr
	paymentConfigMgrLock sync.Mutex
)

// PaymentConfigManager is singleton for PaymentConfigMgr
func PaymentConfigManager() *PaymentConfigMgr {
	if paymentConfigMgr != nil {
		return paymentConfigMgr
	}

	paymentConfigMgrLock.Lock()
	defer paymentConfigMgrLock.Unlock()

	if paymentConfigMgr != nil {
		return paymentConfigMgr
	}

	paymentConfigMgr = &PaymentConfigMgr{}
	return paymentConfigMgr
}

const (
	paymentConfigCollectionName = "payment_config"
)

// Get a PaymentConfig
func (m *PaymentConfigMgr) Get(app int, paymentConfigID string) (paymentConfig *PaymentConfig, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": app, "payment_config_id": paymentConfigID}
	paymentConfig = &PaymentConfig{}
	err = instance.MongoOfficial().Collection(paymentConfigCollectionName).FindOne(ctx, filter).Decode(paymentConfig)
	if err == mongo.ErrNoDocuments {
		paymentConfig = nil
		err = nil
		return
	}
	return
}

// GetWithTx for use in Tx
func (m *PaymentConfigMgr) GetWithTx(sessionContext mongo.SessionContext, app int, paymentConfigID string) (paymentConfig *PaymentConfig, err error) {
	filter := bson.M{"app": app, "payment_config_id": paymentConfigID}
	paymentConfig = &PaymentConfig{}
	err = instance.MongoOfficial().Collection(paymentConfigCollectionName).FindOne(sessionContext, filter).Decode(paymentConfig)
	if err == mongo.ErrNoDocuments {
		paymentConfig = nil
		err = nil
		return
	}
	return
}

// Insert a PaymentConfig
func (m *PaymentConfigMgr) Insert(paymentConfig PaymentConfig) (err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = instance.MongoOfficial().Collection(paymentConfigCollectionName).InsertOne(ctx, paymentConfig)
	return
}

// UpdateWithTx updates a PaymentConfig with Tx
func (m *PaymentConfigMgr) UpdateWithTx(sessionContext mongo.SessionContext, paymentConfig PaymentConfig) (exists bool, err error) {

	filter := bson.M{"app": paymentConfig.App, "payment_config_id": paymentConfig.PaymentConfigID}

	result, err := instance.MongoOfficial().Collection(paymentConfigCollectionName).ReplaceOne(sessionContext, filter, paymentConfig)
	if err != nil {
		return
	}

	exists = result.MatchedCount > 0
	return
}

// DeleteOne delete a PaymentConfig
// only called from test
func (m *PaymentConfigMgr) DeleteOne(app int, paymentConfigID string) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": app, "payment_config_id": paymentConfigID}
	result, err := instance.MongoOfficial().Collection(paymentConfigCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		return
	}

	exists = result.DeletedCount > 0
	return
}

// DeleteOneWithTx delete a PaymentConfig with Tx
func (m *PaymentConfigMgr) DeleteOneWithTx(sessionContext mongo.SessionContext, app int, paymentConfigID string) (exists bool, err error) {
	filter := bson.M{"app": app, "payment_config_id": paymentConfigID}
	result, err := instance.MongoOfficial().Collection(paymentConfigCollectionName).DeleteOne(sessionContext, filter)
	if err != nil {
		return
	}

	exists = result.DeletedCount > 0
	return
}

// Init for this collection
func (m *PaymentConfigMgr) Init() (err error) {

	opts := &options.IndexOptions{}
	opts.SetName("u_app_payment_config_id")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "app", Value: bsonx.Int32(1)}, {Key: "payment_config_id", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	_, err = instance.MongoOfficial().Collection(paymentConfigCollectionName).Indexes().CreateOne(context.Background(), index)
	return
}
