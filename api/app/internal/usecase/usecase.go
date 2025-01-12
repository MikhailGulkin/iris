package usecase

import (
	"context"
	"encoding/json"
)

type EventUC struct {
	Queue  Queue
	Broker Broker
}

func NewEventHandler(queue Queue, broker Broker) *EventUC {
	return &EventUC{
		Queue:  queue,
		Broker: broker,
	}
}

func (h *EventUC) ProcessEvent(ctx context.Context, event Event) error {
	marshal, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return h.Queue.Send(ctx, marshal)
}

func (h *EventUC) ProcessReadMessage(ctx context.Context, msg Message) error {
	return h.Broker.ProducerMessage(ctx, msg)
}
