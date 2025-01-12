package fabric

import (
	"api/app/internal/delivery/websocket"
	"api/app/internal/infra/shared"
	"context"
	"github.com/MikhailGulkin/packages/ws"
)

type PipeProcessorFabric[C Consumer, R ws.ReadPipeProcessor] struct {
	queueSize      int
	consumerFabric ConsumerFabric[C]
	read           ReadProcessorFabric[R]
}

func NewPipeProcessorFabric[C Consumer, R ws.ReadPipeProcessor](
	queueSizeInt int,
	fabric ConsumerFabric[C],
	read ReadProcessorFabric[R],
) *PipeProcessorFabric[C, R] {
	return &PipeProcessorFabric[C, R]{
		queueSize:      queueSizeInt,
		consumerFabric: fabric,
		read:           read,
	}
}

func (p *PipeProcessorFabric[C, R]) NewPipeProcessor(ctx context.Context, uniqueID string) (ws.PipeProcessor, error) {
	queue := shared.NewQueue(p.queueSize)

	consumer, err := p.consumerFabric.CreateConsumer(ctx, queue, uniqueID)
	if err != nil {
		return nil, err
	}
	go consumer.Consume(context.WithoutCancel(ctx))
	readProcessor := p.read.CreateReadProcessor(ctx)
	writeProcessor := websocket.NewWriteProcessor(queue)

	return NewProcessorImpl(
		readProcessor,
		writeProcessor,
		consumer,
	), nil

}

type ProcessorImpl struct {
	r        ws.ReadPipeProcessor
	w        ws.WritePipeProcessor
	consumer Consumer
}

func NewProcessorImpl(
	ReadProcessor ws.ReadPipeProcessor,
	WriteProcessor ws.WritePipeProcessor,
	consumer Consumer,
) *ProcessorImpl {
	return &ProcessorImpl{
		r:        ReadProcessor,
		w:        WriteProcessor,
		consumer: consumer,
	}
}

func (p *ProcessorImpl) ProcessRead(ctx context.Context, messageType int, msg []byte) ([]byte, error) {
	return p.r.ProcessRead(ctx, messageType, msg)
}

func (p *ProcessorImpl) ListenWrite(ctx context.Context) <-chan []byte {
	return p.w.ListenWrite(ctx)
}

func (p *ProcessorImpl) Close() error {
	return p.consumer.Close()
}
