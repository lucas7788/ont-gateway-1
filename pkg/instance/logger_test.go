package instance

import (
	"sync"
	"testing"
)

// go test -race .
func TestLoggerDataRace(t *testing.T) {
	const N= 10
	wg := &sync.WaitGroup{}

	wg.Add(N)
	for i :=0; i< N; i++ {
		go func() {
			defer wg.Done()
			_ = Logger()
		}()
	}

	wg.Wait()
}
