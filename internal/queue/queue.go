package queue

import (
	goqueue "github.com/antonio-alexander/go-queue"
	infinite "github.com/antonio-alexander/go-queue/infinite"
)

func New(queueSize int) interface {
	goqueue.Owner
	goqueue.GarbageCollecter
	goqueue.Dequeuer
	goqueue.Enqueuer
	goqueue.EnqueueInFronter
	goqueue.Length
	goqueue.Event
	goqueue.Peeker
} {
	return infinite.New(queueSize)
}
