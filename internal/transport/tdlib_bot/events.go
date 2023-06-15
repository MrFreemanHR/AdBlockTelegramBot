package tdlibbot

import (
	"adblock_bot/internal/adapter/logger"

	"github.com/zelenin/go-tdlib/client"
)

func (b *tdlibbot) ProcessEvents(event client.Type) {
	if event.GetType() == client.TypeUpdateNewMessage {
		message, ok := (event).(*client.UpdateNewMessage)
		if !ok {
			logger.Logger().Warn("Can't cast new message to struct")
			return
		}
		b.OnMessage(message)
	}
	// logger.Logger().UselessInfo("Data: %#+v\n", event)
}

func (b *tdlibbot) OnMessage(event *client.UpdateNewMessage) {
	logger.Logger().Info("Message: %#+v", event.Message.Content)
}
