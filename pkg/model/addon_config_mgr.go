package model

import (
	"context"
	"sync"

	"github.com/go-redis/redis"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// AddonConfigMgr for AddonConfig crud
type AddonConfigMgr struct {
}

var (
	addonConfigMgr     *AddonConfigMgr
	addonConfigMgrLock sync.Mutex
)

// AddonConfigManager is singleton for AddonConfigMgr
func AddonConfigManager() *AddonConfigMgr {
	if addonConfigMgr != nil {
		return addonConfigMgr
	}

	addonConfigMgrLock.Lock()
	defer addonConfigMgrLock.Unlock()

	if addonConfigMgr != nil {
		return addonConfigMgr
	}

	addonConfigMgr = &AddonConfigMgr{}
	return addonConfigMgr
}

const (
	addonConfigCollectionName = "addon_config"
	cacheTimeout              = 0
)

// Upsert config
func (m *AddonConfigMgr) Upsert(ac AddonConfig) (err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"addon_id": ac.AddonID, "tenant_id": ac.TenantID}

	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "config", Value: ac.Config}},
		}}

	opt := options.Update().SetUpsert(true)

	_, err = instance.MongoOfficial().Collection(addonConfigCollectionName).UpdateOne(ctx, filter, update, opt)
	if err != nil {
		return
	}

	err = instance.RedisCache().Set(keyForAddonConfig(ac.AddonID, ac.TenantID), ac.Config, cacheTimeout)
	return
}

// Delete config
func (m *AddonConfigMgr) Delete(addonID, tenantID string) (err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"addon_id": addonID, "tenant_id": tenantID}
	_, err = instance.MongoOfficial().Collection(addonConfigCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		return
	}

	err = instance.RedisCache().Delete(keyForAddonConfig(addonID, tenantID))
	return
}

func (m *AddonConfigMgr) getNoCache(addonID, tenantID string) (ac *AddonConfig, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ac = &AddonConfig{}
	filter := bson.M{"addon_id": addonID, "tenant_id": tenantID}
	err = instance.MongoOfficial().Collection(addonConfigCollectionName).FindOne(ctx, filter).Decode(ac)
	if err == mongo.ErrNoDocuments {
		ac = nil
		err = nil
		return
	}
	return
}

// Get config
func (m *AddonConfigMgr) Get(addonID, tenantID string) (ac *AddonConfig, err error) {
	key := keyForAddonConfig(addonID, tenantID)
	v, err := instance.RedisCache().Get(key)
	if err != nil {
		if err == redis.Nil {
			ac, err = m.getNoCache(addonID, tenantID)
			if err == nil && ac != nil {
				err = instance.RedisCache().Set(key, ac.Config, cacheTimeout)
			}
			return
		}
		return
	}

	ac = &AddonConfig{AddonID: addonID, TenantID: tenantID, Config: v}

	return
}

// EnsureIndex add index for this collection
func (m *AddonConfigMgr) EnsureIndex() (err error) {

	opts := &options.IndexOptions{}
	opts.SetName("u_addon_id_tenant_id")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "addon_id", Value: bsonx.Int32(1)}, {Key: "tenant_id", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	_, err = instance.MongoOfficial().Collection(addonConfigCollectionName).Indexes().CreateOne(context.Background(), index)
	return
}
