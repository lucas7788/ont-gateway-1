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

// AddonDeploymentMgr for AddonDeployment crud
type AddonDeploymentMgr struct {
}

var (
	addonDeploymentMgr     *AddonDeploymentMgr
	addonDeploymentMgrLock sync.Mutex
)

// AddonDeploymentManager is singleton for AddonDeploymentMgr
func AddonDeploymentManager() *AddonDeploymentMgr {
	if addonDeploymentMgr != nil {
		return addonDeploymentMgr
	}

	addonConfigMgrLock.Lock()
	defer addonConfigMgrLock.Unlock()

	if addonDeploymentMgr != nil {
		return addonDeploymentMgr
	}

	addonDeploymentMgr = &AddonDeploymentMgr{}
	return addonDeploymentMgr
}

const (
	addonDeploymentCollectionName = "addon_deployment"
)

// Upsert deployment
func (m *AddonDeploymentMgr) Upsert(ad AddonDeployment) (err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"addonID": ad.AddonID, "tenantID": ad.TenantID}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "sip", Value: ad.SIP}},
		}}

	opt := options.Update().SetUpsert(true)

	_, err = instance.MongoOfficial().Collection(addonDeploymentCollectionName).UpdateOne(ctx, filter, update, opt)
	if err != nil {
		return
	}

	err = instance.RedisCache().Set(keyForAddonDeployment(ad.AddonID, ad.TenantID), ad.SIP, cacheTimeout)
	return
}

func (m *AddonDeploymentMgr) getNoCache(addonID, tenantID string) (ad *AddonDeployment, exists bool, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ad = &AddonDeployment{}
	filter := bson.M{"addonID": addonID, "tenantID": tenantID}
	err = instance.MongoOfficial().Collection(addonDeploymentCollectionName).FindOne(ctx, filter).Decode(ad)
	if err == mongo.ErrNoDocuments {
		err = nil
		return
	}
	exists = true
	return
}

// Get deployment
func (m *AddonDeploymentMgr) Get(addonID, tenantID string) (ad *AddonDeployment, exists bool, err error) {
	key := keyForAddonDeployment(addonID, tenantID)
	v, err := instance.RedisCache().Get(key)
	if err != nil {
		if err == redis.Nil {
			ad, exists, err = m.getNoCache(addonID, tenantID)
			if exists {
				err = instance.RedisCache().Set(key, ad.SIP, cacheTimeout)
			}
			return
		}
		return
	}

	exists = true
	ad = &AddonDeployment{AddonID: addonID, TenantID: tenantID, SIP: v}

	return
}

// EnsureIndex add index for this collection
func (m *AddonDeploymentMgr) EnsureIndex() (err error) {

	opts := &options.IndexOptions{}
	opts.SetName("i_addon_id_tenant_id")
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "addon_id", Value: bsonx.Int32(1)}, {Key: "tenant_id", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	instance.MongoOfficial().Collection(addonDeploymentCollectionName).Indexes().CreateOne(context.Background(), index)
	return
}
