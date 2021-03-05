package goqueue

import (
	"testing"
)

func TestLen(t *testing.T) {
	q := NewPriorityQueue(3)
	if q.Len() != 0 {
		t.Error("TestLen_len1 <> 0")
	}
	q.Offer(struct{}{}, 1)
	if q.Len() != 1 {
		t.Error("TestLen_len2 <> 1")
	}
}

func TestOffer(t *testing.T) {
	q := NewPriorityQueue(0)
	ele := q.Offer(struct{}{}, 100)
	if ele == nil || ele.Index != 0 {
		t.Error("TestOffer_err1 ele.Index != 0")
	}
	ele = q.Offer(struct{}{}, 2)
	if ele == nil || ele.Index != 0 {
		t.Error("TestOffer_err2 ele.Index != 0", ele)
	}
}

func TestPeek(t *testing.T) {
	q := NewPriorityQueue(3)
	q.Offer(1, 100)
	ele := q.Peek()
	if ele == nil || ele.Value.(int) != 1 {
		t.Error("TestPeek_err1  ele.Value != 1")
	}
	q.Offer(2, 10)
	ele = q.Peek()
	if ele == nil || ele.Value.(int) != 2 {
		t.Error("TestPeek_err2 ele.Value != 2")
	}
}

func TestPeekAndShift(t *testing.T) {
	q := NewPriorityQueue(3)
	ele, wait := q.PeekAndShift(1)
	if ele != nil || wait != 0 {
		t.Errorf("TestPeekAndShift_err1 ele=%+v wait=%d", ele, wait)
	}
	q.Offer(1, 100)
	ele, wait = q.PeekAndShift(1)
	if ele != nil || wait != 99 {
		t.Errorf("TestPeekAndShift_err2 ele=%+v wait=%d", ele, wait)
	}
	ele, wait = q.PeekAndShift(111)
	if ele == nil || wait != 0 {
		t.Errorf("TestPeekAndShift_err3 ele=%+v wait=%d", ele, wait)
	}
}

func TestPeekAndRemove(t *testing.T) {
	q := NewPriorityQueue(3)
	ele := q.PeekAndRemove()
	if ele != nil {
		t.Errorf("TestPeekAndRemove_err1 ele=%+v", ele)
	}
	q.Offer(1, 100)
	q.Offer(2, 99)
	ele = q.PeekAndRemove()
	if ele == nil || ele.Value.(int) != 2 {
		t.Errorf("TestPeekAndRemove_err2 ele=%+v", ele)
	}
	ele = q.PeekAndRemove()
	if ele == nil || ele.Value.(int) != 1 {
		t.Errorf("TestPeekAndRemove_err3 ele=%+v", ele)
	}
	ele = q.PeekAndRemove()
	if ele != nil {
		t.Errorf("TestPeekAndRemove_err4 ele=%+v", ele)
	}
}
