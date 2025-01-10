package scylla

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v3"
	"time"
)

type MessageDAO struct {
	session *gocqlx.Session
}

func NewMessageDAO(session *gocqlx.Session) *MessageDAO {
	return &MessageDAO{session: session}
}

func (m *MessageDAO) GetByID(context context.Context, id uuid.UUID) (Message, error) {
	q := m.session.ContextQuery(context, queryGetMessageByID, nil).Bind(gocql.UUID(id))
	var message Message
	if err := q.GetRelease(&message); err != nil {
		return Message{}, err
	}
	return message, nil
}

func (m *MessageDAO) GetMessages(ctx context.Context, chatID uuid.UUID, cursor *time.Time, limit int) ([]Message, error) {
	query := func() *gocqlx.Queryx {
		if cursor != nil {
			return m.session.ContextQuery(ctx, queryGetMessagesCursor, nil).Bind(
				gocql.UUID(chatID),
				*cursor,
				limit,
			)
		} else {
			return m.session.ContextQuery(ctx, queryGetMessagesByChatID, nil).Bind(
				gocql.UUID(chatID),
				limit,
			)
		}
	}()

	var messages []Message
	if err := query.SelectRelease(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}
