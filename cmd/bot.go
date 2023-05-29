package cmd

import (
	"adblock_bot/infrastructure/sqlite"
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
	db := sqlite.New()

	pingMessageHandler := ping.New(tgbot.GetTelegramAPI())
	cmdCmdMessageHandler := cmdHandler.New(tgbot.GetTelegramAPI(), db)

	messageHandlers := []interfaces.MessageHandler{
		pingMessageHandler,
		cmdCmdMessageHandler,
	}
	tgbot.SetMessageHandlers(messageHandlers)

	tgbot.Run()
}
