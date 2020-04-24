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

// WalletMgr for Wallet crud
type WalletMgr struct {
}

var (
	walletMgr     *WalletMgr
	walletMgrLock sync.Mutex
)

// WalletManager is singleton for WalletMgr
func WalletManager() *WalletMgr {
	if walletMgr != nil {
		return walletMgr
	}

	walletMgrLock.Lock()
	defer walletMgrLock.Unlock()

	if walletMgr != nil {
		return walletMgr
	}

	walletMgr = &WalletMgr{}
	return walletMgr
}

const (
	walletCollectionName = "wallet"
)

// Insert a Wallet
func (m *WalletMgr) Insert(w Wallet) (err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = instance.MongoOfficial().Collection(walletCollectionName).InsertOne(ctx, w)

	return
}

// GetOne returns a Wallet by name
func (m *WalletMgr) GetOne(name string) (w *Wallet, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"name": name}
	w = &Wallet{}
	err = instance.MongoOfficial().Collection(walletCollectionName).FindOne(ctx, filter).Decode(w)
	if err == mongo.ErrNoDocuments {
		w = nil
		err = nil
		return
	}

	return
}

// EnsureIndex add index for this collection
func (m *WalletMgr) EnsureIndex() (err error) {
	hashOpts := &options.IndexOptions{}
	hashOpts.SetName("u-wallet-name")
	hashOpts.SetUnique(true)
	nameUniqueIdx := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "name", Value: bsonx.Int32(1)}},
		Options: hashOpts,
	}

	models := []mongo.IndexModel{nameUniqueIdx}

	_, err = instance.MongoOfficial().Collection(walletCollectionName).Indexes().CreateMany(context.Background(), models)

	return
}
