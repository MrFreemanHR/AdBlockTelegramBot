package cmd

import (
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
	_ = tdlibbot.New()
}
