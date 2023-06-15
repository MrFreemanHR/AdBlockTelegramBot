package bot

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"
	"adblock_bot/internal/core/interfaces"

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

func (b *bot) GetTelegramAPI() *tgbotapi.BotAPI {
	return b.b
}

func (b *bot) Run() {
	for update := range *b.updatesChan {
		b.ProcessEvents(&update)
	}
}
