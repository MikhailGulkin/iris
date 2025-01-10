package aggregate

import (
	"chat/app/internal/domain/entity"
	"fmt"
	"testing"
)

func TestChat_AddMessage(t *testing.T) {
	chat := Chat{}
	message := entity.Message{}
	fmt.Println(chat, message)
}
