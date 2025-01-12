package shared

import (
	"context"
	"errors"
	"testing"
)

func TestQueue_Send(t *testing.T) {
	q := NewQueue(1)

	err := q.Send(context.Background(), []byte("test"))
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	err = q.Send(context.Background(), []byte("test"))
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("expected ErrTimeout, got %v", err)
	}
}

func TestQueue_SendContextCanceled(t *testing.T) {
	q := NewQueue(1)

	ctx, cancel := context.WithCancel(context.Background())

	_ = q.Send(ctx, []byte("test"))
	c := make(chan error)
	go func() {
		c <- q.Send(ctx, []byte("test"))
	}()

	cancel()
	err := <-c

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
