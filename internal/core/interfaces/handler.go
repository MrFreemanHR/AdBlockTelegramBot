package interfaces

import (
	"adblock_bot/internal/core/entity"
)

type MessageHandler interface {
	ProcessMessage(event *entity.TelegramMessage) bool
}
