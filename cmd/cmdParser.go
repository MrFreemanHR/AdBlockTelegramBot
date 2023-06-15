package cmd

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/adapter/parser"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var cmdParserCmd = &cobra.Command{
	Use:     "parser",
	Short:   "Command parser for bot's commands",
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		RunParser()
	},
}

func RunParser() {
	p := parser.New()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		if exit(text) {
			break
		}
		c, err := p.Parse(text)
		if err != nil {
			logger.Logger().Warn("Error: %s", err.Error())
			continue
		}
		printCmdInfo(c)
	}
}

func printCmdInfo(c *parser.Cmd) {
	logger.Logger().Info("=== Command summary ===")
	logger.Logger().Info("Command name: %s", c.Name)
	logger.Logger().Info("=======================")
}

func exit(input string) bool {
	lowerInput := strings.ToLower(input)
	if lowerInput == "q" || lowerInput == "quit" || lowerInput == "exit" {
		return true
	}
	return false
}
