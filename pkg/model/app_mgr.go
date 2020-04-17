package model

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.uber.org/zap"
)

// AppMgr for App crud
type AppMgr struct {
	latest atomic.Value // map[int]App
}

var (
	appMgr     *AppMgr
	appMgrLock sync.Mutex
)

// AppManager is singleton for AppMgr
func AppManager() *AppMgr {
	if appMgr != nil {
		return appMgr
	}

	appMgrLock.Lock()
	defer appMgrLock.Unlock()

	if appMgr != nil {
		return appMgr
	}

	m := &AppMgr{}
	err := m.reload()
	if err != nil {
		logger.Instance().Error("AppMgr.reload", zap.Error(err))
		panic(fmt.Sprintf("AppMgr.reload err:%v", err))
	}

	appMgr = m
	return appMgr
}

const (
	appCollectionName = "app"
)

func (m *AppMgr) updateLatest(apps []App) {
	latest := make(map[int]App)
	for _, app := range apps {
		latest[app.ID] = app
	}
	m.latest.Store(latest)
}

// Insert an App
func (m *AppMgr) Insert(app App) (err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = instance.MongoOfficial().Collection(appCollectionName).InsertOne(ctx, app)
	return
}

func (m *AppMgr) reload() (err error) {
	apps, err := m.GetAllFromDB()
	if err != nil {
		return
	}

	m.updateLatest(apps)
	return
}

// GetAll returns all App from memory
func (m *AppMgr) GetAll() (apps []App) {
	appMap := m.latest.Load().(map[int]App)
	for _, app := range appMap {
		apps = append(apps, app)
	}
	return
}

// GetApp returns App info from memory
func (m *AppMgr) GetApp(id int) (app App, exists bool) {
	appMap := m.latest.Load().(map[int]App)

	app, exists = appMap[id]
	return
}

// GetAllFromDB returns all App from db
func (m *AppMgr) GetAllFromDB() (apps []App, err error) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{}
	cursor, err := instance.MongoOfficial().Collection(appCollectionName).Find(ctx, filter)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &apps)
	return
}

// EnsureIndex add index for this collection
func (m *AppMgr) EnsureIndex() (err error) {

	idOpts := &options.IndexOptions{}
	idOpts.SetName("u_app_id")
	idOpts.SetUnique(true)
	idIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}},
		Options: idOpts,
	}

	models := []mongo.IndexModel{idIndex}

	_, err = instance.MongoOfficial().Collection(appCollectionName).Indexes().CreateMany(context.Background(), models)
	return
}
