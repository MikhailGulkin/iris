package fabric

import (
	"api/app/internal/delivery/websocket"
	"api/app/internal/usecase"
	"context"
	"github.com/MikhailGulkin/packages/ws"
)

type PipeProcessorFabric[C Consumer, Q usecase.Queue, R ws.ReadPipeProcessor] struct {
	queueFabric    QueueFabric[Q]
	consumerFabric ConsumerFabric[C]
	readFabric     ReadProcessorFabric[R]
}

func NewPipeProcessorFabric[C Consumer, Q usecase.Queue, R ws.ReadPipeProcessor](
	queueFabric QueueFabric[Q],
	fabric ConsumerFabric[C],
	read ReadProcessorFabric[R],
) *PipeProcessorFabric[C, Q, R] {
	return &PipeProcessorFabric[C, Q, R]{
		queueFabric:    queueFabric,
		consumerFabric: fabric,
		readFabric:     read,
	}
}

func (p *PipeProcessorFabric[C, Q, R]) NewPipeProcessor(ctx context.Context, uniqueID string) (ws.PipeProcessor, error) {
	queue := p.queueFabric.CreateQueue(ctx)

	consumer, err := p.consumerFabric.CreateConsumer(ctx, queue, uniqueID)
	if err != nil {
		return nil, err
	}
	go consumer.Consume(context.WithoutCancel(ctx))
	readProcessor := p.readFabric.CreateReadProcessor(ctx)
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
