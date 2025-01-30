package usecase

import (
	"context"
	"encoding/json"
)

type EventUC struct {
	Queue Queue
}

func NewEventHandler(queue Queue) *EventUC {
	return &EventUC{
		Queue: queue,
	}
}

func (h *EventUC) ProcessEvent(ctx context.Context, event Event) error {
	marshal, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return h.Queue.Send(ctx, marshal)
}

type MessageUC struct {
	Broker Broker
}

func NewMessageHandler(broker Broker) *MessageUC {
	return &MessageUC{
		Broker: broker,
	}
}

func (m *MessageUC) ProcessReadMessage(ctx context.Context, msg Message) error {
	return m.Broker.ProducerMessage(ctx, msg)
}
