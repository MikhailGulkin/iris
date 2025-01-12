package main

import (
	"api/app/cmd/fabric"
	kafkaP "api/app/internal/infra/kafka"
	"context"
	"github.com/MikhailGulkin/packages/kafka/producer"
	log "github.com/MikhailGulkin/packages/logger"
	"github.com/MikhailGulkin/packages/rabbit"
	"github.com/MikhailGulkin/packages/ws"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.Default()
	valid := validator.New()
	kafkaWriter, err := producer.NewProducer(producer.Config{
		Brokers: []string{"localhost:9092"},
		Topic:   "user.messages",
	})
	if err != nil {
		return
	}
	kafkaProducer := kafkaP.NewProducer(kafkaWriter)
	pipeProcessorFabric := fabric.NewPipeProcessorFabric(
		rabbit.Config{
			URL:          "amqp://guest:guest@localhost:5672/",
			Exchange:     "user.messages",
			QueuePattern: "user.id",
		},
		100,
		logger,
		valid,
		kafkaProducer,
	)

	manager := ws.NewManager(ws.WithProcessorFabric(pipeProcessorFabric))
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
		err := manager.Process(uuid.New().String(), c.Writer, c.Request, nil)
		if err != nil {
			logger.Errorw("error processing ws", "error", err)
			return
		}
		return
	})
	go func() {
		app.Run(":8000")
	}()

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh
}
