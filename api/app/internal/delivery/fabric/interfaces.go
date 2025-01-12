package fabric

import (
	"api/app/internal/usecase"
	"context"
)

type ConsumerFabric[C any] interface {
	CreateConsumer(ctx context.Context, queue usecase.Queue, uniqueID string) (C, error)
}
type Consumer interface {
	Consume(ctx context.Context)
	Close() error
}
type ReadProcessorFabric[R any] interface {
	CreateReadProcessor(ctx context.Context) R
}
