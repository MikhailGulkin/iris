package scylla

import (
	"github.com/gocql/gocql"
	"time"
)

type Message struct {
	ID          gocql.UUID
	ChatID      gocql.UUID
	UserID      int64
	EditAt      *time.Time
	DeletedAt   *time.Time
	Text        string
	ContentType string
	SendAt      time.Time
	CreatedAt   time.Time
}
