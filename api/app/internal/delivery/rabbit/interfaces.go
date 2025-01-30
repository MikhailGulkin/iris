package rabbit

import (
	"api/app/internal/usecase"
	"context"
)

type EventHandler interface {
	ProcessEvent(ctx context.Context, event usecase.Event) error
}
