package cmd

import (
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/transport/bot/message_handler/ping"
	tdlibbot "adblock_bot/internal/transport/tdlib_bot"

	"github.com/spf13/cobra"
)

var userBot = &cobra.Command{
	Use:     "userbot",
	Short:   "Start user-bot with authentication by phone number",
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		StartUserBot()
	},
}

func StartUserBot() {
	userbot := tdlibbot.New()

	if userbot == nil {
		return
	}

	pingMessageHandler := ping.New(userbot.GetTelegramAPI())

	messageHandlers := []interfaces.MessageHandler{
		pingMessageHandler,
	}
	userbot.SetMessageHandlers(messageHandlers)

	userbot.Run()
}
