package kafka

import (
	"api/app/internal/usecase"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func NewProducer(w *kafka.Writer) *Producer {
	return &Producer{w: w}
}

func (p *Producer) ProducerMessage(ctx context.Context, msg usecase.Message) error {
	bytes, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	return p.w.WriteMessages(ctx, kafka.Message{
		Value: bytes,
	})
}
