package rabbit

import (
	"api/app/internal/usecase"
	"context"
	"encoding/json"
	log "github.com/MikhailGulkin/packages/logger"
	"github.com/MikhailGulkin/packages/rabbit"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerHandler struct {
	logger    log.Logger
	validator *validator.Validate
	handler   EventHandler
}

func NewConsumerHandler(
	logger log.Logger,
	validator *validator.Validate,
	handler EventHandler,
) *ConsumerHandler {
	return &ConsumerHandler{
		logger:    logger,
		validator: validator,
		handler:   handler,
	}
}

func (c *ConsumerHandler) Consume(ctx context.Context, msg <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msg:
			if !ok {
				c.logger.Info("channel closed")
				return
			}
			func() {
				var event usecase.Event
				err := json.Unmarshal(msg.Body, &event)
				if err != nil {
					c.logger.Errorw("failed to unmarshal event", "error", err)
					return
				}
				if err := c.validator.Struct(event); err != nil {
					c.logger.Errorw("failed to validate event", "error", err)
					return
				}
				err = c.handler.ProcessEvent(ctx, event)
				if err != nil {
					c.logger.Errorw("failed to handle event", "error", err)
					return
				}
				err = msg.Ack(false)
				if err != nil {
					c.logger.Errorw("error ack message", "error", err)
				}
			}()
		}
	}
}

type ConsumerHandlerFabric struct {
	rabbitConfig rabbit.Config
	logger       log.Logger
	validator    *validator.Validate
}

type ConsumerHandlerFacade struct {
	handler *ConsumerHandler
	conn    *rabbit.Conn
	msg     <-chan amqp.Delivery
}

func (c *ConsumerHandlerFacade) Consume(ctx context.Context) {
	c.handler.Consume(ctx, c.msg)
}

func (c *ConsumerHandlerFacade) Close() error {
	return c.conn.Close()
}

func (c *ConsumerHandlerFabric) CreateConsumer(
	_ context.Context,
	queue usecase.Queue,
	uniqueID string,
) (*ConsumerHandlerFacade, error) {
	conn, err := rabbit.NewRabbitCh(c.rabbitConfig)
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
	rabbitHandler := usecase.NewEventHandler(queue)

	handler := NewConsumerHandler(c.logger, c.validator, rabbitHandler)

	return &ConsumerHandlerFacade{
		msg:     msg,
		conn:    conn,
		handler: handler,
	}, nil
}

func NewConsumerHandlerFabric(
	rabbitConfig rabbit.Config,
	logger log.Logger,
	validator *validator.Validate,
) *ConsumerHandlerFabric {
	return &ConsumerHandlerFabric{
		rabbitConfig: rabbitConfig,
		logger:       logger,
		validator:    validator,
	}
}
