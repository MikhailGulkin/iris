package main

import (
	log "api/app/pkg/logger"
	"api/app/pkg/ws"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				p.send <- []byte("ping")
			}
		}
	}()

	return p.send
}

type PipeProcessorFabricImpl struct{}

func (p *PipeProcessorFabricImpl) NewPipeProcessor(ctx context.Context, userID int) (ws.PipeProcessor, error) {
	return NewProcessorImpl(), nil
}

func main() {
	logger := log.Default()

	manager := ws.NewManager(
		ws.WithProcessorFabric(&PipeProcessorFabricImpl{}),
	)
	defer func() {
		err := manager.Close()
		if err != nil {
			logger.Errorw("error closing manager", "error", err)
		}
		logger.Infow("manager closed")
	}()
	app := gin.Default()
	go func() {
		manager.Run(context.Background())
	}()
	app.GET("/ws", func(c *gin.Context) {
		err := manager.Process(rand.Int(), c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("error", err)
		}
	})
	go func() {
		app.Run(":8000")
	}()

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh

}
