package goqueue

import (
	"sync"
	"testing"
	"time"
)

func TestDelayQueueOffer(t *testing.T) {
	q := NewDelayQueue(3)
	go q.Poll()
	if len(q.C) != 0 {
		t.Error("TestOffer1 len <> 0")
	}
	q.Offer(1, time.Now())
	ele := <-q.C
	if ele == nil {
		t.Error("TestOffer2 len <> 0")
	}
}

func TestStop(t *testing.T) {
	q := NewDelayQueue(3)
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan struct{})
	go func() {
		q.Poll()
		ch <- struct{}{}
	}()
	go func() {
		wg.Wait()
		q.Stop()
	}()
	wg.Done()
	select {
	case <-time.After(10 * time.Second):
		t.Error("TestStop_err  ")
	case <-ch:
	}
}

func TestPoll(t *testing.T) {
	return
}
