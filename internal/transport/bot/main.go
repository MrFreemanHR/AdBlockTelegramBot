package bot

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"
	"adblock_bot/internal/core/entity"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/transport"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bot struct {
	b               *tgbotapi.BotAPI
	updatesChan     *tgbotapi.UpdatesChannel
	messageHandlers []interfaces.MessageHandler
}

func New() *bot {
	tgbot, err := tgbotapi.NewBotAPI(config.CurrentConfig.Token)
	if err != nil {
		logger.Logger().Fatal("Can't connect to Telegram: %s", err.Error())
		return nil
	}

	logger.Logger().Info("Logged in: %s", tgbot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updatesChan := tgbot.GetUpdatesChan(u)

	return &bot{
		b:           tgbot,
		updatesChan: &updatesChan,
	}
}

func (b *bot) SetMessageHandlers(messageHandlers []interfaces.MessageHandler) {
	b.messageHandlers = messageHandlers
}

func (b *bot) GetTelegramAPI() *transport.TelegramAPI {
	self, err := b.b.GetMe()
	if err != nil {
		logger.Logger().Fatal("[TG API] Can't get me: %s", err.Error())
		return nil
	}
	return transport.NewUnifiedTransportWithBotAPI(
		b.b,
		entity.TelegramUser{
			ID:           self.ID,
			IsBot:        self.IsBot,
			FirstName:    self.FirstName,
			LastName:     self.LastName,
			UserName:     self.UserName,
			LanguageCode: self.LanguageCode,
		},
	)
}

func (b *bot) Run() {
	for update := range *b.updatesChan {
		b.ProcessEvents(&update)
	}
}
