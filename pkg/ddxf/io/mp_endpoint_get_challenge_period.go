package io

import (
	"time"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
)

// MPEndpointGetChallengePeriodInput ...
type MPEndpointGetChallengePeriodInput struct {
}

// MPEndpointGetChallengePeriodOutput ...
type MPEndpointGetChallengePeriodOutput struct {
	Period time.Duration  `bson:"period" json:"period"`
	io.BaseResp
}
