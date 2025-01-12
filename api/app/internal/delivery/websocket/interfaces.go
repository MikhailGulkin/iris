package websocket

import (
	"api/app/internal/usecase"
	"context"
)

type MessageReceiver interface {
	Receive(ctx context.Context) <-chan []byte
}

type ReadProcessorHandler interface {
	ProcessReadMessage(ctx context.Context, msg usecase.Message) error
}
