package misc

import (
	"testing"

	"gotest.tools/assert"
)

func TestChain(t *testing.T) {

	amount, err := NewOntSdk().GetAmountTransferred("ff250fdf93f4e888b1f1f94af9792fad84d64406fe2ae409204f844b6876b1f1")
	assert.Assert(t, err == nil && amount == 10000000)

}
