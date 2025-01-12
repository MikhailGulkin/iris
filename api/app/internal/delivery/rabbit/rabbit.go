package rabbit

import (
	"api/app/internal/usecase"
	"context"
	"encoding/json"
	log "github.com/MikhailGulkin/packages/logger"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerHandler struct {
	msg       <-chan amqp.Delivery
	logger    log.Logger
	validator *validator.Validate
	handler   EventHandler
}

func NewConsumerHandler(msg <-chan amqp.Delivery, logger log.Logger, validator *validator.Validate, handler EventHandler) *ConsumerHandler {
	return &ConsumerHandler{
		msg:       msg,
		logger:    logger,
		validator: validator,
		handler:   handler,
	}
}

func (c *ConsumerHandler) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.msg:
			if !ok {
				c.logger.Error("channel closed unexpectedly")
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
