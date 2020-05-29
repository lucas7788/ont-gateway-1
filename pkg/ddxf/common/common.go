package common

import "github.com/satori/go.uuid"

func GenerateUUId() string {
	return uuid.NewV4().String()
}
