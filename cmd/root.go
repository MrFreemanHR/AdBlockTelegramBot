package cmd

import "github.com/spf13/cobra"

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
}

func initConfig() {

}
