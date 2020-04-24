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

	filter := bson.M{"addon_id": ad.AddonID, "tenant_id": ad.TenantID}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "sip", Value: ad.SIP},
			bson.E{Key: "updated_at", Value: ad.UpdatedAt},
		}}}

	opt := options.Update().SetUpsert(true)

	_, err = instance.MongoOfficial().Collection(addonDeploymentCollectionName).UpdateOne(ctx, filter, update, opt)
	if err != nil {
		return
	}

	err = instance.RedisCache().Set(keyForAddonDeployment(ad.AddonID, ad.TenantID), ad.SIP, cacheTimeout)
	return
}

func (m *AddonDeploymentMgr) getNoCache(addonID, tenantID string) (ad *AddonDeployment, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ad = &AddonDeployment{}
	filter := bson.M{"addon_id": addonID, "tenant_id": tenantID}
	err = instance.MongoOfficial().Collection(addonDeploymentCollectionName).FindOne(ctx, filter).Decode(ad)
	if err == mongo.ErrNoDocuments {
		ad = nil
		err = nil
		return
	}
	return
}

// Get deployment
func (m *AddonDeploymentMgr) Get(addonID, tenantID string) (ad *AddonDeployment, exists bool, err error) {
	key := keyForAddonDeployment(addonID, tenantID)
	v, err := instance.RedisCache().Get(key)
	if err != nil {
		if err == redis.Nil {
			ad, err = m.getNoCache(addonID, tenantID)
			if ad != nil {
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

// Init for this collection
func (m *AddonDeploymentMgr) Init() (err error) {

	opts := &options.IndexOptions{}
	opts.SetName("u_addon_id_tenant_id")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "addon_id", Value: bsonx.Int32(1)}, {Key: "tenant_id", Value: bsonx.Int32(1)}},
		Options: opts,
	}

	_, err = instance.MongoOfficial().Collection(addonDeploymentCollectionName).Indexes().CreateOne(context.Background(), index)
	return
}
