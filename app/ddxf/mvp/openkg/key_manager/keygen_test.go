package key_manager

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func TestGetKeyPair(t *testing.T) {
	for i := 0; i < 10000; i++ {
		in := rand.Int()
		pri, pub := GetKeyPair([]byte(strconv.Itoa(in)))
		if pri == nil || pub == nil {
			panic("here")
		}
	}
}

func Test_getKeyPair(t *testing.T) {
	pri, _ := GetKeyPair([]byte("123"))
	for i := 0; i < 100000; i++ {
		pri2, _ := GetKeyPair([]byte("123"))
		assert.Equal(t, pri, pri2)
	}
}
