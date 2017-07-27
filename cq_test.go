package cq

import (
	"math/rand"
	"testing"
	"time"
)

func TestEnqueue(t *testing.T) {
	q := New()

	if q.Len() != 0 {
		t.Error("Length != 0")
	}
	items := q.ToSlice()
	if len(items) != 0 {
		t.Error("0 items expected, found:", len(items))
	}

	q.Enqueue(1)
	q.Enqueue(2)

	if q.Len() != 2 {
		t.Error("Length != 2")
	}
	items = q.ToSlice()
	if len(items) != 2 {
		t.Error("2 items expected, found:", len(items))
	}
}

func TestDequeue(t *testing.T) {
	q := New()
	q.Enqueue(1)
	q.Enqueue(2)

	a, ok := q.Dequeue()
	if a != 1 || !ok {
		t.Error("Expected a: 1 ok: true. found a:", a, "ok:", ok)
	}

	b, ok := q.Dequeue()
	if b != 2 || !ok {
		t.Error("Expected b: 2 ok: true. found b:", b, "ok:", ok)
	}

	c, ok := q.Dequeue()
	if c != nil || ok {
		t.Error("Expected c: nil ok: false. found c:", c, "ok:", ok)
	}

	if q.Len() != 0 {
		t.Error("Length != 0")
	}
	items := q.ToSlice()
	if len(items) != 0 {
		t.Error("0 items expected, found:", len(items))
	}
}

func TestConcurrentEnqueue(t *testing.T) {
	q := New()
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go enqueueRoutine(q, i, done)
	}
	for i := 0; i < 10; i++ {
		<-done
	}

	if q.Len() != 10 {
		t.Error("Length != 10")
	}

	items := q.ToSlice()
	if len(items) != 10 {
		t.Error("10 items expected, found:", len(items))
	}
}

func TestConcurrentDequeue(t *testing.T) {
	q := New()
	res := make(chan interface{}, 15)
	for i := 0; i < 10; i++ {
		q.Enqueue(i)
	}

	for i := 0; i < 5; i++ {
		go dequeueRoutine(q, res)
	}

	resArr := make([]interface{}, 0)
	for i := 0; i < 15; i++ {
		resArr = append(resArr, <-res)
	}

	if q.Len() != 0 {
		t.Error("Length != 0")
	}
	items := q.ToSlice()
	if len(items) != 0 {
		t.Error("0 items expected, found:", len(items))
	}

	itemCount := 0
	nilCount := 0
	for i := 0; i < 15; i++ {
		item := resArr[i]

		if _, ok := item.(int); ok {
			itemCount += 1
		} else {
			nilCount += 1
		}
	}
	if itemCount != 10 || nilCount != 5 {
		t.Error("Expected number of items: 10, found:", itemCount)
		t.Error("Expected number of nils: 5, found:", nilCount)
	}
}

func dequeueRoutine(q *Queue, res chan interface{}) {
	t := rand.Int63() % 10
	time.Sleep(time.Duration(t) * time.Millisecond)
	for i := 0; i < 3; i++ {
		item, _ := q.Dequeue()
		res <- item
	}
}

func enqueueRoutine(q *Queue, num int, done chan bool) {
	t := rand.Int63() % 10
	time.Sleep(time.Duration(t) * time.Millisecond)
	q.Enqueue(num)
	done <- true
}

func BenchmarkEnqueue(b *testing.B) {
	q := New()
	var msg [100]byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(msg)
	}
	b.Log(len(q.ToSlice()))
}

func BenchmarkEnqueueParallel(b *testing.B) {
	q := New()
	var msg [100]byte
	f := func(pb *testing.PB) {
		for pb.Next() {
			q.Enqueue(msg)
		}
	}
	b.RunParallel(f)
	b.Log(len(q.ToSlice()))
}
