package model

import (
	"testing"

	"gotest.tools/assert"
)

func TestPaymentManager(t *testing.T) {
	_, err := PaymentManager().QueryToNotifyPreRecharging(10)
	assert.Assert(t, err == nil, err)
}
