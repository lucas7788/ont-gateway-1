package key_manager

import (
	"github.com/ontio/go-bip32"
)

func GenerateSeed() ([]byte, error) {
	return bip32.NewSeed()
}

func ClearSeed(seed []byte) {
	size := len(seed)
	for i := 0; i < size; i++ {
		seed[i] = 0
	}
}

func ClearKey(key *bip32.Key) {
	size := len(key.Key)
	for i := 0; i < size; i++ {
		key.Key[i] = 0
	}
}
