package crypto

import (
	"crypto/sha256"
)

// SHA256 returns sha256 result for bytes
func SHA256(bytes []byte) []byte {
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}
