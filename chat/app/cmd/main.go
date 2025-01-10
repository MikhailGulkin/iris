package main

import (
	scylla2 "chat/app/internal/infra/scylla"
	"chat/app/pkg/scylla"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func main() {
	//fmt.Println(personTable.Get())
	//return
	client, err := scylla.NewScyllaClient(scylla.Config{
		Hosts: []string{"localhost:9042"},
	})
	if err != nil {
		panic(err)
	}
	cursor, err := time.Parse("2006-01-02 15:04:05.000", "2025-01-07 22:47:50.880")
	if err != nil {
		panic(err)
	}

	daoMessage := scylla2.NewMessageDAO(client)
	fmt.Println(daoMessage.GetMessages(
		context.Background(),
		uuid.MustParse("3df36dac-973e-450a-8d2b-b561511a888c"),
		&cursor,
		10,
	))
	fmt.Println(daoMessage.GetByID(
		context.Background(),
		uuid.MustParse("3df26dac-973e-450a-8d2b-b561511a888c"),
	))
}
