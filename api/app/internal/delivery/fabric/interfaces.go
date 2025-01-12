package fabric

import (
	"api/app/internal/usecase"
	"context"
	"github.com/MikhailGulkin/packages/ws"
)

type ConsumerFabric[C Consumer] interface {
	CreateConsumer(ctx context.Context, queue usecase.Queue, uniqueID string) (C, error)
}
type Consumer interface {
	Consume(ctx context.Context)
	Close() error
}
type ReadProcessorFabric[R ws.ReadPipeProcessor] interface {
	CreateReadProcessor(ctx context.Context) R
}
type QueueFabric[Q usecase.Queue] interface {
	CreateQueue(ctx context.Context) Q
}
