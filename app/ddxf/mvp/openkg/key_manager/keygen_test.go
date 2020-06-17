package key_manager

import (
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

}
