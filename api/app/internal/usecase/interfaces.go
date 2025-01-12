package usecase

import "context"

type Queue interface {
	Send(ctx context.Context, msg []byte) error
}

type Broker interface {
	ProducerMessage(ctx context.Context, msg Message) error
}
