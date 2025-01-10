package aggregate

import (
	"chat/app/internal/domain/entity"
	"github.com/google/uuid"
	"time"
)

type Chat struct {
	ID              uuid.UUID
	Name            string
	ConversationID  uuid.UUID
	Messages        []entity.Message
	Type            ChatType
	PinnedMessageID *uuid.UUID
	DeletedAt       *time.Time
}

func (c *Chat) AddMessages(messages ...entity.Message) {
	c.Messages = append(c.Messages, messages...)
}

func (c *Chat) PinMessage(messageID uuid.UUID) {
	c.PinnedMessageID = &messageID
}

func (c *Chat) DeleteMessage(messageID uuid.UUID) {
	for i, m := range c.Messages {
		if m.ID == messageID {
			c.Messages[i].Delete()
		}
	}
}
