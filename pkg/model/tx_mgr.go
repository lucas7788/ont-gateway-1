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

// TxMgr for Tx crud
type TxMgr struct {
}

var (
	txMgr     *TxMgr
	txMgrLock sync.Mutex
)

// TxManager is singleton for TxMgr
func TxManager() *TxMgr {
	if txMgr != nil {
		return txMgr
	}

	txMgrLock.Lock()
	defer txMgrLock.Unlock()

	if txMgr != nil {
		return txMgr
	}

	txMgr = &TxMgr{}
	return txMgr
}

const (
	txCollectionName = "tx"
)

// Upsert a Tx
func (m *TxMgr) Upsert(tx Tx) (err error) {

	_, err = m.update(tx, true)

	return
}

// Update a Tx
func (m *TxMgr) Update(tx Tx) (exists bool, err error) {

	result, err := m.update(tx, false)
	if err != nil {
		return
	}

	exists = result.ModifiedCount > 0
	return
}

// UpdateNotifyError only updates NotifyErrxxx
func (m *TxMgr) UpdateNotifyError(txHash string, msg string) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"hash": txHash}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "notify_err_msg", Value: msg},
			bson.E{Key: "updated_at", Value: time.Now()},
		}},
		bson.E{Key: "$inc", Value: bson.D{
			bson.E{Key: "notify_err_count", Value: 1},
		}}}

	updateResult, err := instance.MongoOfficial().Collection(txCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return
	}

	exists = updateResult.ModifiedCount > 0
	return
}

// UpdatePollError only updates PollErrxxx
func (m *TxMgr) UpdatePollError(txHash string, msg string) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"hash": txHash}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "poll_err_msg", Value: msg},
			bson.E{Key: "updated_at", Value: time.Now()},
		}},
		bson.E{Key: "$inc", Value: bson.D{
			bson.E{Key: "poll_err_count", Value: 1},
		}}}

	updateResult, err := instance.MongoOfficial().Collection(txCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return
	}

	exists = updateResult.ModifiedCount > 0
	return
}

// FinishPoll is called when poll finish, whether succeed or fail
func (m *TxMgr) FinishPoll(txHash string, result TxPollResult, errMsg string) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"hash": txHash}
	setV := bson.D{
		bson.E{Key: "state", Value: TxStateToNotify},
		bson.E{Key: "result", Value: result},
		bson.E{Key: "updated_at", Value: time.Now()},
	}
	var update bson.D
	if errMsg != "" {
		setV = append(setV, bson.E{Key: "poll_err_msg", Value: errMsg})
		update = bson.D{
			bson.E{Key: "$set", Value: setV},
			bson.E{Key: "$inc", Value: bson.D{
				bson.E{Key: "poll_err_count", Value: 1},
			}}}
	} else {
		update = bson.D{
			bson.E{Key: "$set", Value: setV}}
	}

	updateResult, err := instance.MongoOfficial().Collection(txCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return
	}

	exists = updateResult.ModifiedCount > 0
	return
}

// UpdateState only updates TxState (and updated_at)
func (m *TxMgr) UpdateState(txHash string, state TxState) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"hash": txHash}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "state", Value: state},
			bson.E{Key: "updated_at", Value: time.Now()},
		}}}

	updateResult, err := instance.MongoOfficial().Collection(txCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return
	}

	exists = updateResult.ModifiedCount > 0
	return
}

// UpdateResult only updates TxPollResult (and updated_at)
func (m *TxMgr) UpdateResult(txHash string, result TxPollResult) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"hash": txHash}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "result", Value: result},
			bson.E{Key: "updated_at", Value: time.Now()},
		}}}

	updateResult, err := instance.MongoOfficial().Collection(txCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return
	}

	exists = updateResult.ModifiedCount > 0
	return
}

func (m *TxMgr) update(tx Tx, upsert bool) (result *mongo.UpdateResult, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"hash": tx.Hash}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "hash", Value: tx.Hash},
			bson.E{Key: "state", Value: tx.State},
			bson.E{Key: "expire_at", Value: tx.ExpireAt},
			bson.E{Key: "updated_at", Value: tx.UpdatedAt},
		}}}

	if upsert {
		opt := options.Update().SetUpsert(true)

		result, err = instance.MongoOfficial().Collection(txCollectionName).UpdateOne(ctx, filter, update, opt)
	} else {
		result, err = instance.MongoOfficial().Collection(txCollectionName).UpdateOne(ctx, filter, update)
	}

	return
}

// QueryToPoll for query topoll tx
func (m *TxMgr) QueryToPoll(n int) (txlist []Tx, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"state": TxStateToPoll}
	options := options.Find()
	// sort by poll_err_count asc, updated_at asc
	options.SetSort(bson.D{bson.E{Key: "poll_err_count", Value: 1}, bson.E{Key: "updated_at", Value: 1}})
	cursor, err := instance.MongoOfficial().Collection(txCollectionName).Find(ctx, filter, options)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &txlist)
	return
}

// QueryToNotify for query tonotify tx
func (m *TxMgr) QueryToNotify(n int) (txlist []Tx, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"state": TxStateToNotify}
	options := options.Find()
	// sort by notify_err_count asc, updated_at asc
	options.SetSort(bson.D{bson.E{Key: "notify_err_count", Value: 1}, bson.E{Key: "updated_at", Value: 1}})
	cursor, err := instance.MongoOfficial().Collection(txCollectionName).Find(ctx, filter, options)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &txlist)
	return
}

// Delete by txHash
func (m *TxMgr) Delete(txHash string) (exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"hash": txHash}
	result, err := instance.MongoOfficial().Collection(txCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		return
	}

	exists = result.DeletedCount > 0
	return
}

// EnsureIndex add index for this collection
func (m *TxMgr) EnsureIndex() (err error) {

	hashOpts := &options.IndexOptions{}
	hashOpts.SetName("u-tx-hash")
	hashOpts.SetUnique(true)
	hashIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "hash", Value: bsonx.Int32(1)}},
		Options: hashOpts,
	}

	// used by QueryToPoll
	stateOpts := &options.IndexOptions{}
	stateOpts.SetName("i-tx-state-poll_err_count-updated_at")
	stateAndPollErrCountAndUpdatedAtIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "state", Value: bsonx.Int32(1)}, {Key: "poll_err_count", Value: bsonx.Int32(1)}, {Key: "updated_at", Value: bsonx.Int32(1)}},
		Options: stateOpts,
	}

	// used by QueryToNotify
	resultOpts := &options.IndexOptions{}
	resultOpts.SetName("i-tx-state-notify_err_count-updated_at")
	stateAndResultAndNotifyErrCountAndUpdatedAtIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "state", Value: bsonx.Int32(1)}, {Key: "notify_err_count", Value: bsonx.Int32(1)}, {Key: "updated_at", Value: bsonx.Int32(1)}},
		Options: resultOpts,
	}

	models := []mongo.IndexModel{hashIndex, stateAndPollErrCountAndUpdatedAtIndex, stateAndResultAndNotifyErrCountAndUpdatedAtIndex}

	_, err = instance.MongoOfficial().Collection(txCollectionName).Indexes().CreateMany(context.Background(), models)
	return
}
