package ping

import (
	"adblock_bot/internal/adapter/logger"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler struct {
	bot *tgbotapi.BotAPI
}

func New(bot *tgbotapi.BotAPI) *handler {
	return &handler{
		bot: bot,
	}
}

func (h *handler) ProcessMessage(event *tgbotapi.Update) bool {
	messageText := event.Message.Text
	if strings.Contains(messageText, h.bot.Self.UserName) &&
		strings.Contains(strings.ToLower(messageText), "ping") {
		reply := tgbotapi.NewMessage(event.Message.Chat.ID, "Pong")
		reply.ReplyToMessageID = event.Message.MessageID
		_, err := h.bot.Send(reply)
		if err != nil {
			logger.Logger().Warn("Can't send message: %s", err)
		}
	}
	return false
}
