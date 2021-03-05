package goqueue

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func demo() {
	type Task struct {
		ID   int
		Time time.Time
	}

	// create a delay queue with cap 3
	q := NewDelayQueue(3)
	// add a Task to queue which should run after 10 seconds later

	var wg sync.WaitGroup
	wg.Add(100)
	// run a goroutine to deal timer check
	go q.Poll()
	go func() {
		for ele := range q.C {
			if ele != nil {
				fmt.Println(time.Now(), ele.(Task).Time)
				wg.Done()
			}
		}
	}()
	for i := 0; i < 100; i++ {
		r := rand.Intn(100)
		t := Task{
			ID:   r,
			Time: time.Now().Add(time.Duration(r) * time.Second),
		}
		q.Offer(t, t.Time)
	}
	wg.Wait()
}
