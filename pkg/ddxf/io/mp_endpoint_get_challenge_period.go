package io

import (
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"time"
)

// MPEndpointGetChallengePeriodInput ...
type MPEndpointGetChallengePeriodInput struct {
}

// MPEndpointGetChallengePeriodOutput ...
type MPEndpointGetChallengePeriodOutput struct {
	Period time.Duration `bson:"period" json:"period"`
	io.BaseResp
}
