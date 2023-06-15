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
	var subcmd = c.ChildSubcmd
	var lastSubcmd *parser.Subcmd
	for {
		logger.Logger().Info("Subcmd: %s", subcmd.Name)
		if subcmd.Child == nil {
			lastSubcmd = subcmd
			break
		}
		subcmd = subcmd.Child
	}
	logger.Logger().Info("Required args:")
	for i, requiredArg := range lastSubcmd.RequiredArgs {
		logger.Logger().Info("[%d] Value: %s", i, requiredArg.GetValue())
	}
	logger.Logger().Info("Optional args:")
	for i, optionalArg := range lastSubcmd.OptionalArgs {
		logger.Logger().Info("[%d] Value: %s", i, optionalArg.GetValue())
	}
	logger.Logger().Info("=======================")
}

func exit(input string) bool {
	lowerInput := strings.ToLower(input)
	if lowerInput == "q" || lowerInput == "quit" || lowerInput == "exit" {
		return true
	}
	return false
}
