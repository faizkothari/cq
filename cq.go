package cq

import (
	"sync/atomic"
	"unsafe"
)

type (
	node struct {
		data interface{}
		next unsafe.Pointer
	}

	Queue struct {
		head, tail unsafe.Pointer
		len        int64
	}
)

func New() *Queue {
	q := &Queue{}
	q.tail = newNode(nil)
	q.head = q.tail
	return q
}

func newNode(data interface{}) unsafe.Pointer {
	return unsafe.Pointer(&node{data: data, next: unsafe.Pointer(nil)})
}

func (q *Queue) Enqueue(data interface{}) {
	if data == nil {
		return
	}
	n := newNode(data)

	for {
		tail := q.tail
		next := (*node)(tail).next
		if tail == q.tail {
			if next == unsafe.Pointer(nil) {
				if atomic.CompareAndSwapPointer(&(*node)(q.tail).next, next, unsafe.Pointer(n)) {
					atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(n))
					atomic.AddInt64(&q.len, 1)
					break
				}
			} else {
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}
		}
	}
}

func (q *Queue) Dequeue() interface{} {
	for {
		head := q.head
		tail := q.tail
		next := (*node)(head).next

		if head == q.head {
			if head == tail {
				if next == unsafe.Pointer(nil) {
					return nil
				}
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			} else {
				data := (*node)(next).data
				if atomic.CompareAndSwapPointer(&q.head, head, next) {
					atomic.AddInt64(&q.len, -1)
					return data
				}
			}
		}
	}
}

func (q *Queue) Len() int64 {
	return atomic.LoadInt64(&q.len)
}

func (q *Queue) ToSlice() []interface{} {
	s := make([]interface{}, 0, 4)
	p := (*node)(atomic.LoadPointer(&q.head)).next
	for ; p != unsafe.Pointer(nil); p = (*node)(p).next {
		s = append(s, (*node)(p).data)
	}
	return s
}
