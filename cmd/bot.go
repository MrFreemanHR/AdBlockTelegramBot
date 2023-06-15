package cmd

import (
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/transport/bot"
	"adblock_bot/internal/transport/bot/message_handler/ping"

	"github.com/spf13/cobra"
)

var botCmd = &cobra.Command{
	Use:     "bot",
	Short:   "Start bot",
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		StartBot()
	},
}

func StartBot() {
	tgbot := bot.New()

	pingMessageHandler := ping.New(tgbot.GetTelegramAPI())

	messageHandlers := []interfaces.MessageHandler{
		pingMessageHandler,
	}
	tgbot.SetMessageHandlers(messageHandlers)

	tgbot.Run()
}
