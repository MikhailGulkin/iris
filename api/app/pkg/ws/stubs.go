package ws

import (
	"context"
	"fmt"
)

type ProcessorImpl struct {
	send chan []byte
}

func NewProcessorImpl() *ProcessorImpl {
	return &ProcessorImpl{
		send: make(chan []byte, 256),
	}
}

func (p *ProcessorImpl) ProcessRead(ctx context.Context, messageType int, msg []byte) error {
	fmt.Println("ProcessRead", messageType, msg)
	return nil
}

func (p *ProcessorImpl) ProcessWrite() <-chan []byte {
	return p.send
}

type PipeProcessorFabricImpl struct{}

func (p *PipeProcessorFabricImpl) NewPipeProcessor(ctx context.Context, userID int) (PipeProcessor, error) {
	return NewProcessorImpl(), nil
}
