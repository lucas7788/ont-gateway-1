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

// PaymentOrderMgr for PaymentOrder crud
type PaymentOrderMgr struct {
}

var (
	paymentOrderMgr     *PaymentOrderMgr
	paymentOrderMgrLock sync.Mutex
)

// PaymentOrderManager is singleton for PaymentOrderMgr
func PaymentOrderManager() *PaymentOrderMgr {
	if paymentOrderMgr != nil {
		return paymentOrderMgr
	}

	paymentOrderMgrLock.Lock()
	defer paymentOrderMgrLock.Unlock()

	if paymentOrderMgr != nil {
		return paymentOrderMgr
	}

	paymentOrderMgr = &PaymentOrderMgr{}
	return paymentOrderMgr
}

const (
	paymentOrderCollectionName = "payment_order"
)

// InsertWithTx inserts a PaymentOrder within Tx
func (m *PaymentOrderMgr) InsertWithTx(sessionContext mongo.SessionContext, order PaymentOrder) (err error) {
	_, err = instance.MongoOfficial().Collection(paymentOrderCollectionName).InsertOne(sessionContext, order)
	return
}

// GetPaymentOrders returns all orders for paymentID
func (m *PaymentOrderMgr) GetPaymentOrders(sessionContext mongo.SessionContext, app int, paymentID string) (orders []PaymentOrder, err error) {
	filter := bson.M{"app": app, "payment_id": paymentID}
	opts := options.Find()
	// Sort by `_id` field descending
	opts.SetSort(bson.D{bson.E{Key: "_id", Value: -1}})
	cursor, err := instance.MongoOfficial().Collection(paymentOrderCollectionName).Find(sessionContext, filter, opts)
	if err != nil {
		return
	}

	err = cursor.All(sessionContext, &orders)

	return
}

// TotalAmountPaidWithTx returns the total paid amount
func (m *PaymentOrderMgr) TotalAmountPaidWithTx(sessionContext mongo.SessionContext, app int, paymentID string) (amount int, err error) {
	filter := bson.M{"app": app, "payment_id": paymentID}
	cursor, err := instance.MongoOfficial().Collection(paymentOrderCollectionName).Find(sessionContext, filter)
	if err != nil {
		return
	}

	var orders []PaymentOrder
	err = cursor.All(sessionContext, &orders)
	if err != nil {
		return
	}

	for _, order := range orders {
		amount += order.Amount
	}
	return
}

// DeletePaymentOrders for delete payment orders by paymentID
// only called from test
func (m *PaymentOrderMgr) DeletePaymentOrders(app int, paymentID string) (n int64, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": app, "payment_id": paymentID}

	result, err := instance.MongoOfficial().Collection(paymentOrderCollectionName).DeleteMany(ctx, filter)
	if err != nil {
		return
	}

	n = result.DeletedCount
	return
}

// Init for this collection
func (m *PaymentOrderMgr) Init() (err error) {

	opts := &options.IndexOptions{}
	opts.SetName("u_app_order_id")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "app", Value: bsonx.Int32(1)}, {Key: "order_id", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	_, err = instance.MongoOfficial().Collection(paymentOrderCollectionName).Indexes().CreateOne(context.Background(), index)
	return
}
