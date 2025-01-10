package ws

import "context"

type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

type PipeProcessor interface {
	ProcessRead(ctx context.Context, messageType int, msg []byte) error
	ProcessWrite() <-chan []byte
}

type PipeProcessorFabric interface {
	NewPipeProcessor(ctx context.Context, userID int) (PipeProcessor, error)
}

type Client interface {
	GetClientID() string
	Run(ctx context.Context) error
	Close() error
}
