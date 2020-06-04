package param

import (
	"testing"

	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
)

func TestTokenTemplate_Serialize(t *testing.T) {

}

func TestTokenResourceTyEndpoint_Serialize(t *testing.T) {

}

func TestResourceDDO_Serialize(t *testing.T) {
	ddo := &ResourceDDO{
		Manager:                  common.ADDRESS_EMPTY,
		TokenResourceTyEndpoints: []*TokenResourceTyEndpoint{},
		ItemMetaHash:             common.UINT256_EMPTY,
		DTC:                      common.ADDRESS_EMPTY,
		MP:                       common.ADDRESS_EMPTY,
		Split:                    common.ADDRESS_EMPTY,
	}
	sink := common.NewZeroCopySink(nil)
	ddo.Serialize(sink)

	ddo2 := &ResourceDDO{}
	source := common.NewZeroCopySource(sink.Bytes())
	err := ddo2.Deserialize(source)
	assert.Nil(t, err)
	assert.Equal(t, ddo, ddo2)
}

func TestFee_Serialize(t *testing.T) {

}

func TestDTokenItem_Serialize(t *testing.T) {

}
