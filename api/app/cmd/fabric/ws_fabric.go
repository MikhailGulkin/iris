package fabric

import (
	rabbitD "api/app/internal/delivery/rabbit"
	"api/app/internal/delivery/websocket"
	"api/app/internal/infra/shared"
	"api/app/internal/usecase"
	"context"
	log "github.com/MikhailGulkin/packages/logger"
	"github.com/MikhailGulkin/packages/rabbit"
	"github.com/MikhailGulkin/packages/ws"
	"github.com/go-playground/validator/v10"
)

type PipeProcessorFabric struct {
	rabbitConfig rabbit.Config
	queueSize    int
	logger       log.Logger
	validator    *validator.Validate
	broker       usecase.Broker
}

func NewPipeProcessorFabric(
	rabbitConfig rabbit.Config,
	queueSize int,
	logger log.Logger,
	validator *validator.Validate,
	broker usecase.Broker,
) *PipeProcessorFabric {
	return &PipeProcessorFabric{
		rabbitConfig: rabbitConfig,
		queueSize:    queueSize,
		logger:       logger,
		validator:    validator,
		broker:       broker,
	}
}

func (p *PipeProcessorFabric) NewPipeProcessor(ctx context.Context, uniqueID string) (ws.PipeProcessor, error) {
	conn, err := rabbit.NewRabbitCh(p.rabbitConfig)
	if err != nil {
		return nil, err
	}
	err = conn.DeclareAndBindQueue(uniqueID, "")
	if err != nil {
		return nil, err
	}
	msg, err := conn.Consume("", "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	queue := shared.NewQueue(p.queueSize)
	rabbitHandler := usecase.NewEventHandler(queue, p.broker)

	socketConsumer := rabbitD.NewConsumerHandler(
		msg,
		p.logger,
		p.validator,
		rabbitHandler,
	)
	go socketConsumer.Consume(context.WithoutCancel(ctx))

	readProcessor := websocket.NewReadProcessor(rabbitHandler, p.validator)
	writeProcessor := websocket.NewWriteProcessor(queue)

	return NewProcessorImpl(
		readProcessor,
		writeProcessor,
		conn,
	), nil

}

type ProcessorImpl struct {
	ws.ReadPipeProcessor
	ws.WritePipeProcessor
	conn     *rabbit.Conn
	consumer *rabbitD.ConsumerHandler
}

func NewProcessorImpl(
	ReadProcessor ws.ReadPipeProcessor,
	WriteProcessor ws.WritePipeProcessor,
	conn *rabbit.Conn,
) *ProcessorImpl {
	return &ProcessorImpl{
		ReadPipeProcessor:  ReadProcessor,
		WritePipeProcessor: WriteProcessor,
		conn:               conn,
	}
}

func (p *ProcessorImpl) Close() error {
	return p.conn.Close()
}
