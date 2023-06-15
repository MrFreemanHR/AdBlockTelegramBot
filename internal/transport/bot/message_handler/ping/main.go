package ping

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/transport"
	"strings"

	"adblock_bot/internal/core/entity"
)

type handler struct {
	api *transport.TelegramAPI
}

func New(api *transport.TelegramAPI) interfaces.MessageHandler {
	return &handler{
		api: api,
	}
}

func (h *handler) ProcessMessage(event *entity.TelegramMessage) bool {
	messageText := event.Text
	if strings.Contains(messageText, h.api.Self().UserName) &&
		strings.Contains(strings.ToLower(messageText), "ping") {
		reply := entity.NewMessage(event.Chat.ID, "Pong")
		reply.ReplyToMessage = event
		err := h.api.SendMessage(reply)
		if err != nil {
			logger.Logger().Warn("Can't send message: %s", err)
		}
		return true
	}
	return false
}
