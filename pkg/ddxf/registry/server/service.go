package server

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"context"
	"encoding/hex"
	"net/http"
	"github.com/ontio/ontology/core/signature"
)

const (
	registryCollectionName = "registry"
)


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

	err = signature.Verify(pk, []byte(input.MP), input.Sign)
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
	return
}

