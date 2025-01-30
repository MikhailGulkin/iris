package websocket

import "context"

type WriteProcessor struct {
	receiver MessageReceiver
}

func NewWriteProcessor(receiver MessageReceiver) *WriteProcessor {
	return &WriteProcessor{
		receiver: receiver,
	}
}

func (w *WriteProcessor) ListenWrite(ctx context.Context) <-chan []byte {
	return w.receiver.Receive(ctx)
}
