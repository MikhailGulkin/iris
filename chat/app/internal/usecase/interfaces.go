package usecase

import "chat/app/internal/domain/aggregate"

type ChatRepository interface {
	AddChat(chat aggregate.Chat) error
}
