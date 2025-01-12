package websocket

import (
	"api/app/internal/usecase"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type ReadProcessor struct {
	handler   ReadProcessorHandler
	validator *validator.Validate
}

func NewReadProcessor(
	handler ReadProcessorHandler,
	validator *validator.Validate,
) *ReadProcessor {
	return &ReadProcessor{
		handler:   handler,
		validator: validator,
	}
}

func (r *ReadProcessor) ProcessRead(ctx context.Context, messageType int, msg []byte) ([]byte, error) {
	switch messageType {
	case websocket.TextMessage:
		var message usecase.Message
		err := jsoniter.Unmarshal(msg, &message)
		if err != nil {
			return nil, err
		}
		if err := r.validator.Struct(message); err != nil {
			return nil, err
		}

		err = r.handler.ProcessReadMessage(ctx, message)
		if err != nil {
			return nil, err
		}
		return []byte("ok"), nil
	}

	return nil, ErrInvalidMessage
}
