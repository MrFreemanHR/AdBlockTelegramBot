package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *bot) ProcessEvents(event *tgbotapi.Update) {
	if event.Message != nil {
		b.OnMessage(event)
	}
}

func (b *bot) OnMessage(event *tgbotapi.Update) {
	for _, handler := range b.messageHandlers {
		processed := handler.ProcessMessage(event)
		if processed {
			break
		}
	}
}
