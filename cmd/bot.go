package cmd

import (
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/transport/bot"
	cmdHandler "adblock_bot/internal/transport/bot/message_handler/cmd"
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
	localesCmdMessageHandler := cmdHandler.New(tgbot.GetTelegramAPI())

	messageHandlers := []interfaces.MessageHandler{
		pingMessageHandler,
		localesCmdMessageHandler,
	}
	tgbot.SetMessageHandlers(messageHandlers)

	tgbot.Run()
}
