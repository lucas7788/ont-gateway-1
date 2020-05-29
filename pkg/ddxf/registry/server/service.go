package server

import (
	"context"
	"encoding/hex"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/core/signature"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
)

const (
	registryCollectionName = "registry"
)

// Init for this collection
func Init() (err error) {
	opts := &options.IndexOptions{}
	opts.SetName("u-mp")
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "mp", Value: bsonx.Int32(1)}},
		Options: opts,
	}
	_, err = instance.MongoOfficial().Collection(registryCollectionName).Indexes().CreateOne(context.Background(), index)
	return
}

func AddEndpointService(input io.RegistryAddEndpointInput) (output io.RegistryAddEndpointOutput) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	pkBs, err := hex.DecodeString(input.PubKey)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	_, err = keypair.DeserializePublicKey(pkBs)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	_, err = instance.MongoOfficial().Collection(registryCollectionName).InsertOne(ctx, input)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func RemoveEndpointService(input io.RegistryRemoveEndpointInput) (output io.RegistryRemoveEndpointOutput) {
	timeout := config.Load().MongoConfig.Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	input2 := io.MPRemoveRegistryInput(input)
	ad := &io.MPAddRegistryInput{}
	filter := bson.M{"mp": input2.MP}
	err := instance.MongoOfficial().Collection(registryCollectionName).FindOne(ctx, filter).Decode(ad)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	pkBs, err := hex.DecodeString(ad.PubKey)
	if err != nil {
		panic(err)
	}
	pk, err := keypair.DeserializePublicKey(pkBs)
	if err != nil {
		panic(err)
	}

	sig, err := hex.DecodeString(input.Sign)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	err = signature.Verify(pk, []byte(input.MP), sig)
	if err != nil {
		output.Code = http.StatusBadRequest
		output.Msg = err.Error()
		return
	}
	filter = bson.M{"mp": input.MP}
	_, err = instance.MongoOfficial().Collection(registryCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
	}
	return
}

func QueryEndpointsService(input io.RegistryQueryEndpointsInput) (output io.RegistryQueryEndpointsOutput) {
	filter := bson.D{}
	opts := options.Find()
	// Sort by `_id` field descending
	opts.SetSort(bson.D{bson.E{Key: "_id", Value: -1}})
	cursor, err := instance.MongoOfficial().Collection(registryCollectionName).Find(context.Background(), filter, opts)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	itemMetas := make([]io.MPEndpoint, 0)
	err = cursor.All(context.Background(), &itemMetas)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	output.Endpoints = itemMetas
	return
}
