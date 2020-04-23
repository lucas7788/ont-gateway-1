package model

import (
	"context"
	"sync"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// PaymentMgr for Payment crud
type PaymentMgr struct {
}

var (
	paymentMgr     *PaymentMgr
	paymentMgrLock sync.Mutex
)

// PaymentManager is singleton for PaymentMgr
func PaymentManager() *PaymentMgr {
	if paymentMgr != nil {
		return paymentMgr
	}

	paymentMgrLock.Lock()
	defer paymentMgrLock.Unlock()

	if paymentMgr != nil {
		return paymentMgr
	}

	paymentMgr = &PaymentMgr{}
	return paymentMgr
}

const (
	paymentCollectionName = "payment"
)

// CountForPaymentConfigWithTx returns number of payments for paymentConfigID
func (m *PaymentMgr) CountForPaymentConfigWithTx(sessionContext mongo.SessionContext, app int, paymentConfigID string) (n int64, err error) {

	filter := bson.M{"app": app, "payment_config_id": paymentConfigID}
	n, err = instance.MongoOfficial().Collection(paymentCollectionName).CountDocuments(sessionContext, filter)
	return

}

// Insert a Payment
func (m *PaymentMgr) Insert(sessionContext mongo.SessionContext, payment Payment) (err error) {
	_, err = instance.MongoOfficial().Collection(paymentCollectionName).InsertOne(sessionContext, payment)
	return
}

// UpdateBalanceAndStartTime for update balance and start_time and balance_expire_time
func (m *PaymentMgr) UpdateBalanceAndStartTime(sessionContext mongo.SessionContext, app int, paymentID string, balance int, startTime, balanceExpireTime time.Time) (exists bool, err error) {
	filter := bson.M{"app": app, "payment_id": paymentID}

	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "state", Value: PaymentStateStarted},
			bson.E{Key: "balance", Value: balance},
			bson.E{Key: "start_time", Value: startTime},
			bson.E{Key: "updated_at", Value: time.Now()},
			bson.E{Key: "balance_expire_time", Value: balanceExpireTime}},
		}}

	result, err := instance.MongoOfficial().Collection(paymentCollectionName).UpdateOne(sessionContext, filter, update)
	if err != nil {
		return
	}

	exists = result.MatchedCount > 0
	return
}

// UpdateBalanceAndExpireTime for update balance and balance_expire_time
func (m *PaymentMgr) UpdateBalanceAndExpireTime(sessionContext mongo.SessionContext, app int, paymentID string, balance int, balanceExpireTime time.Time) (exists bool, err error) {
	filter := bson.M{"app": app, "payment_id": paymentID}

	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "balance", Value: balance},
			bson.E{Key: "balance_expire_time", Value: balanceExpireTime},
			bson.E{Key: "updated_at", Value: time.Now()},
			bson.E{Key: "notify_date", Value: ""}},
		}}

	result, err := instance.MongoOfficial().Collection(paymentCollectionName).UpdateOne(sessionContext, filter, update)
	if err != nil {
		return
	}

	exists = result.MatchedCount > 0
	return
}

// UpdateBalance for update balance
func (m *PaymentMgr) UpdateBalance(sessionContext mongo.SessionContext, app int, paymentID string, balance int) (exists bool, err error) {
	filter := bson.M{"app": app, "payment_id": paymentID}

	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "balance", Value: balance},
			bson.E{Key: "updated_at", Value: time.Now()},
		}}}

	result, err := instance.MongoOfficial().Collection(paymentCollectionName).UpdateOne(sessionContext, filter, update)
	if err != nil {
		return
	}

	exists = result.MatchedCount > 0
	return
}

// GetOnePayment finds one Payment by paymentID
func (m *PaymentMgr) GetOnePayment(sessionContext mongo.SessionContext, app int, paymentID string) (payment *Payment, err error) {
	filter := bson.M{"app": app, "payment_id": paymentID}
	payment = &Payment{}
	err = instance.MongoOfficial().Collection(paymentCollectionName).FindOne(sessionContext, filter).Decode(payment)
	if err == mongo.ErrNoDocuments {
		payment = nil
		err = nil
		return
	}
	if err != nil {
		return
	}
	return
}

const (
	notifyDuration = time.Hour * 24 * 30
)

// QueryBalanceExpired returns Payment with expired BalanceExpireTime
func (m *PaymentMgr) QueryBalanceExpired(batch int64) (payments []Payment, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{
		"state": PaymentStateStarted,
		"balance_expire_time": bson.M{
			"$lt": time.Now(),
		},
	}
	opts := options.Find()
	opts.SetLimit(batch)
	cursor, err := instance.MongoOfficial().Collection(paymentCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &payments)
	if err != nil {
		return
	}
	return
}

// QueryToNotifyRecharging returns Payment to notify recharging when balance is negative
func (m *PaymentMgr) QueryToNotifyRecharging(batch int64) (payments []Payment, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{
		"state":       PaymentStateStarted,
		"balance":     bson.M{"$lt": 0},
		"notify_date": bson.M{"$ne": time.Now().Format("2006-01-02")},
	}
	opts := options.Find()
	// sort by notify_err_count asc
	opts.SetSort(bson.D{bson.E{Key: "notify_err_count", Value: 1}})
	opts.SetLimit(batch)
	cursor, err := instance.MongoOfficial().Collection(paymentCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &payments)
	if err != nil {
		return
	}
	return
}

// UpdateNotifyDate for update notify date
func (m *PaymentMgr) UpdateNotifyDate(app int, paymentID string) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": app, "payment_id": paymentID}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "notify_date", Value: nowDate()},
			bson.E{Key: "updated_at", Value: time.Now()},
		}},
		bson.E{Key: "$inc", Value: bson.D{
			bson.E{Key: "notify_err_count", Value: 1},
		}}}

	updateResult, err := instance.MongoOfficial().Collection(paymentCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return
	}

	exists = updateResult.ModifiedCount > 0
	return
}

// UpdateNotifyError for update notify error
func (m *PaymentMgr) UpdateNotifyError(app int, paymentID, errMsg string) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": app, "payment_id": paymentID}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "notify_err_msg", Value: errMsg},
			bson.E{Key: "updated_at", Value: time.Now()},
		}},
		bson.E{Key: "$inc", Value: bson.D{
			bson.E{Key: "notify_err_count", Value: 1},
		}}}

	updateResult, err := instance.MongoOfficial().Collection(paymentCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return
	}

	exists = updateResult.ModifiedCount > 0
	return
}

func nowDate() string {
	return time.Now().Format("2006-01-02")
}

// QueryToNotifyPreRecharging returns Payment to notify pre recharging
func (m *PaymentMgr) QueryToNotifyPreRecharging(batch int64) (payments []Payment, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	now := time.Now()
	filter := bson.M{
		"state":      PaymentStateStarted,
		"pay_period": bson.M{"$gte": PayPeriodMonthly},
		"balance_expire_time": bson.M{
			"$lt": now.Add(notifyDuration),
			"$gt": now,
		},
		"balance": bson.M{"$gte": 0},
		"$expr": bson.M{
			"$gt": []string{"$unit_amount", "$balance"},
		},
		"notify_date": bson.M{"$ne": nowDate()},
	}
	opts := options.Find()
	// sort by notify_err_count asc
	opts.SetSort(bson.D{bson.E{Key: "notify_err_count", Value: 1}})
	opts.SetLimit(batch)
	cursor, err := instance.MongoOfficial().Collection(paymentCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &payments)
	if err != nil {
		return
	}
	return
}

// DeleteOne for delete a payment by paymentID
// only called from test
func (m *PaymentMgr) DeleteOne(app int, paymentID string) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"app": app, "payment_id": paymentID}

	result, err := instance.MongoOfficial().Collection(paymentCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		return
	}

	exists = result.DeletedCount > 0
	return
}

// EnsureIndex add index for this collection
func (m *PaymentMgr) EnsureIndex() (err error) {

	opts := &options.IndexOptions{}
	opts.SetName("u_app_payment_id")
	opts.SetUnique(true)
	appPaymentIDIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "app", Value: bsonx.Int32(1)}, {Key: "payment_id", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	opts = &options.IndexOptions{}
	opts.SetName("i-state-balance_expire_time-balance-unit_amount")
	stateBalanceExpireBalanceUnitAmountIndex := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "state", Value: bsonx.Int32(1)},
			{Key: "balance_expire_time", Value: bsonx.Int32(1)},
			{Key: "balance", Value: bsonx.Int32(1)},
			{Key: "unit_amount", Value: bsonx.Int32(1)},
		},
		Options: opts,
	}

	models := []mongo.IndexModel{
		appPaymentIDIndex,
		stateBalanceExpireBalanceUnitAmountIndex,
	}
	_, err = instance.MongoOfficial().Collection(paymentCollectionName).Indexes().CreateMany(context.Background(), models)
	return
}
