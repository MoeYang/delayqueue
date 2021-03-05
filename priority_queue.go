//Package goqueue 优先级队列(小顶堆)
package goqueue

import (
	"container/heap"
)

// Element
type Element struct {
	Value    interface{}
	Priority int64
	Index    int
}

// PriorityQueue no concurrency safty, need to Lock when use with multi-goroutines
type PriorityQueue []*Element

// NewPriorityQueue NewPriorityQueue
func NewPriorityQueue(size int) PriorityQueue {
	return make(PriorityQueue, 0, size)
}

// Offer add a val with priority
func (pq *PriorityQueue) Offer(val interface{}, priority int64) *Element {
	ele := &Element{Value: val, Priority: priority}
	heap.Push(pq, ele)
	return ele
}

// Peek return a *Element which has the least Priority in queue
func (pq PriorityQueue) Peek() *Element {
	if pq.Len() == 0 {
		return nil
	}
	return pq[0]
}

// PeekAndShift return a *Element which has the least Priority in queue and ele.Priority not more than param priority.
// Otherwise, return nil and priority-ele.Priority
func (pq *PriorityQueue) PeekAndShift(priority int64) (*Element, int64) {
	if pq.Len() == 0 {
		return nil, 0
	}
	ele := (*pq)[0]
	if ele.Priority <= priority {
		heap.Remove(pq, 0)
		return ele, 0
	}
	return nil, ele.Priority - priority
}

// PeekAndShift return the first *Element and delete it from queue
func (pq *PriorityQueue) PeekAndRemove() *Element {
	if pq.Len() == 0 {
		return nil
	}
	ele := (*pq)[0]
	heap.Remove(pq, 0)
	return ele
}

// implements func for container/heap interface

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	// change ele index
	pq[i].Index = j
	pq[j].Index = i
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(ele interface{}) {
	// set index
	l := len(*pq)
	eleP := ele.(*Element)
	eleP.Index = l
	*pq = append(*pq, eleP)
}

func (pq *PriorityQueue) Pop() interface{} {
	l := len(*pq)
	c := cap(*pq)
	// reduce cap
	if l < c/2 && c > 25 {
		tmp := make(PriorityQueue, l, c/2)
		copy(tmp, *pq)
		*pq = tmp
	}
	ele := (*pq)[l-1]
	ele.Index = -1
	*pq = (*pq)[:l-1]
	return ele
}
