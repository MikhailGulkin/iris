package fabric

import (
	"api/app/internal/infra/shared"
	"context"
)

type QueueFabricImpl struct {
	queueSize int
}

func NewQueueFabricImpl(queueSize int) *QueueFabricImpl {
	return &QueueFabricImpl{
		queueSize: queueSize,
	}
}

func (q *QueueFabricImpl) CreateQueue(ctx context.Context) *shared.Queue {
	return shared.NewQueue(q.queueSize)
}
