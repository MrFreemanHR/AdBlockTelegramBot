package cmd

import (
	"adblock_bot/infrastructure/mysql"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/core/messageHandlers/cmd"
	"adblock_bot/internal/core/messageHandlers/ping"
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
	db := mysql.New()

	if userbot == nil {
		return
	}

	pingMessageHandler := ping.New(userbot.GetTelegramAPI())
	cmdMessageHandler := cmd.New(userbot.GetTelegramAPI(), db)

	messageHandlers := []interfaces.MessageHandler{
		pingMessageHandler,
		cmdMessageHandler,
	}
	userbot.SetMessageHandlers(messageHandlers)

	userbot.Run()
}
