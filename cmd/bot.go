package cmd

import (
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/core/messageHandlers/ping"
	"adblock_bot/internal/transport/bot"

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
	// db := mysql.New()

	pingMessageHandler := ping.New(tgbot.GetTelegramAPI())
	// cmdCmdMessageHandler := cmdHandler.New(tgbot.GetTelegramAPI(), db)
	// verifierMessageHandler := verifier.New(tgbot.GetTelegramAPI(), db)

	messageHandlers := []interfaces.MessageHandler{
		// verifierMessageHandler,
		pingMessageHandler,
		// cmdCmdMessageHandler,
	}
	tgbot.SetMessageHandlers(messageHandlers)

	tgbot.Run()
}
