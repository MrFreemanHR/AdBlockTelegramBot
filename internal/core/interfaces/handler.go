package interfaces

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageHandler interface {
	ProcessMessage(event *tgbotapi.Update) bool
}
