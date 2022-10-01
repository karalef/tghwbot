package queue

import (
	"sync"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	var wg sync.WaitGroup
	q := New(func(int) {
		wg.Done()
	}, 10)
	wg.Add(1000)
	go func() {
		for i := 0; i < 1000; i++ {
			q.Push(i)
		}
	}()
	wg.Wait()

	c := make(chan struct{})
	go func() {
		defer close(c)

	}()
	select {
	case <-c:
	case <-time.After(2 * time.Second):
		t.Fail()
	}
}
