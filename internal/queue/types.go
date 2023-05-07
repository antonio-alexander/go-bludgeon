package queue

import goqueue "github.com/antonio-alexander/go-queue"

type Queue interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Peeker
	goqueue.Length
}
