package cq

import (
	_ "fmt"
	"math/rand"
	"testing"
	"time"
)

func TestEnqueue(t *testing.T) {
	q := New()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Enqueue(4)

	//	t.Error(q.ToSlice(), q.Len())
}

func TestDequeue(t *testing.T) {
	q := New()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Enqueue(4)

	//a := q.Dequeue()
	//b := q.Dequeue()
	//	t.Error(q.ToSlice(), q.Len(), a, b)
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
	//	//t.Error(q.ToSlice(), q.Len())
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
	//	t.Error(resArr, q.ToSlice(), q.Len())
}

func dequeueRoutine(q *Queue, res chan interface{}) {
	//t := rand.Int63() % 100
	//time.Sleep(time.Duration(t) * time.Millisecond)
	res <- q.Dequeue()
	res <- q.Dequeue()
	res <- q.Dequeue()
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
