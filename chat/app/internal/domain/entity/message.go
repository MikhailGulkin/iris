package entity

import (
	"chat/app/internal/domain/value_object"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	Content   value_object.Content
	UserID    int64
	EditAt    *time.Time
	DeletedAt *time.Time
	IsNew     bool
}

func CreateNewMessage(chatID uuid.UUID, userID int64, content value_object.Content) Message {
	return NewMessage(uuid.New(), chatID, userID, content, true, nil, nil)
}

func NewMessage(
	id uuid.UUID,
	chatID uuid.UUID,
	userID int64,
	content value_object.Content,
	isNew bool,
	editAt *time.Time,
	deletedAt *time.Time,
) Message {
	return Message{
		ID:        id,
		ChatID:    chatID,
		UserID:    userID,
		Content:   content,
		EditAt:    editAt,
		DeletedAt: deletedAt,
		IsNew:     isNew,
	}
}

func (m *Message) Edit(content value_object.Content) {
	m.Content = content
	now := time.Now()
	m.EditAt = &now
}

func (m *Message) Delete() {
	now := time.Now()
	m.DeletedAt = &now
}
