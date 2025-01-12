package shared

import (
	"context"
	"time"
)

type Queue struct {
	ch chan []byte
}

func NewQueue(queueSize int) *Queue {
	return &Queue{
		ch: make(chan []byte, queueSize),
	}
}

func (q *Queue) Send(ctx context.Context, msg []byte) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case q.ch <- msg:
			return nil
		case <-time.After(5 * time.Second):
			return ErrTimeout
		}
	}
}

func (q *Queue) Receive(ctx context.Context) <-chan []byte {
	return q.ch
}
