package cmd

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "adblock-tg",
	Short:   "Telegram bot for removing ad's from chats",
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func StartApp() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(botCmd)
}

func initConfig() {
	var err error
	config.CurrentConfig, err = config.ParseConfig("config.json")
	cobra.CheckErr(err)
	logger.New(logger.VerbosityLevel(config.CurrentConfig.VerbosityLevel))
}
