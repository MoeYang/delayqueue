//Package goqueue 延迟队列
package goqueue

import (
	"sync"
	"sync/atomic"
	"time"
)

// DelayQueue make your element delay to be done with a Millisecond precision
type DelayQueue struct {
	lock        sync.Mutex
	pq          PriorityQueue
	sleepStatus int32
	wakeUpChan  chan struct{}
	stopChan    chan struct{}

	C chan interface{} // you will get the ready element from this chan
}

// NewDelayQueue New DelayQueue with size
func NewDelayQueue(size int) *DelayQueue {
	return &DelayQueue{
		pq:         NewPriorityQueue(size),
		wakeUpChan: make(chan struct{}),
		stopChan:   make(chan struct{}, 1), // only need one signal, and need Stop() func unblock
		C:          make(chan interface{}),
	}
}

// Offer add a element to queue, runTime is the future time we want to get the element
func (dq *DelayQueue) Offer(value interface{}, runTime time.Time) {
	// lock to add ele
	dq.lock.Lock()
	ele := dq.pq.Offer(value, time2Millisecond(runTime))
	idx := ele.Index
	dq.lock.Unlock()
	// if ele is the first of priQueue, wakeup the sleeping process to check whether the first ele need to be exec.
	// Only if Poll func is sleeping, should send a wakeUpChan signal.
	if idx == 0 {
		// if and only if set sleeping to 0, should send signal
		if atomic.CompareAndSwapInt32(&dq.sleepStatus, 1, 0) {
			dq.wakeUpChan <- struct{}{}
		}
	}
}

// Stop stop to check delayQueue, break loop in Poll func
func (dq *DelayQueue) Stop() {
	select {
	case dq.stopChan <- struct{}{}:
	default:
	}
}

// Poll send the ready element to a chan,
// you should call this func with a new goroutine to make process unblock.
func (dq *DelayQueue) Poll() {
	for {
		now := time2Millisecond(time.Now())
		dq.lock.Lock()
		// get the first ele if ele`s exectime <= nowTime, then remove the ele from queue
		ele, waitTime := dq.pq.PeekAndShift(now)
		if ele == nil {
			// if no ele ready, need to sleep
			atomic.StoreInt32(&dq.sleepStatus, 1)
		}
		dq.lock.Unlock()
		// 1、ele == nil, waitTime == 0 : no ele in queue, need to sleep
		// 2、ele == nil, waitTime > 0 : need to sleep waitTime ms
		// 3、ele != nil, get a ele which need to exec
		if ele == nil {
			if waitTime == 0 {
				select {
				case <-dq.wakeUpChan:
					continue
				case <-dq.stopChan:
					goto exit
				}
			} else {
				select {
				case <-time.After(time.Duration(waitTime) * time.Millisecond):
					// maybe there is a go call Offer() and blocking on dq.wakeUpChan<-struct{}
					// so if someone change dq.sleepStatus to 0, need get the sig from wakeUpChan to unblock the caller
					if atomic.LoadInt32(&dq.sleepStatus) == 0 {
						<-dq.wakeUpChan
					}
					continue
				case <-dq.wakeUpChan:
					continue
				case <-dq.stopChan:
					goto exit
				}
			}
		}
		// has ready ele, send to C
		select {
		case dq.C <- ele.Value:
		case <-dq.stopChan:
			goto exit
		}
	}
exit:
	// if exit, set sleepStatus to 0
	atomic.StoreInt32(&dq.sleepStatus, 0)
}

// time2Millisecond translate time to Millisecond stamp
func time2Millisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
