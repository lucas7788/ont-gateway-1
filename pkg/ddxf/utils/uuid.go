package utils

import (
	"github.com/satori/go.uuid"
)

const (
	SELLER_PUBLISH_PRIFIX = "seller_publish"
)

const (
	UUID_TYPE_RAW             int32 = 1
	UUID_TOKEN_SELLER_PUBLISH int32 = 2
)

func GenerateUUId(uuidType int32) string {
	u1 := uuid.NewV4()
	switch uuidType {
	case UUID_TYPE_RAW:
		return u1.String()
	case UUID_TOKEN_SELLER_PUBLISH:
		return SELLER_PUBLISH_PRIFIX + u1.String()
	}

	return u1.String()
}
