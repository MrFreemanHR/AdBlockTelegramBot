package localeshandler

import (
	"adblock_bot/internal/adapter/locales"
	"adblock_bot/internal/adapter/parser"
	"adblock_bot/internal/core/interfaces"
)

type cmdhandler struct {
}

func New() interfaces.CmdHandler {
	return &cmdhandler{}
}

func (h *cmdhandler) ProcessCommand(cmd *parser.Cmd) string {
	keySubCmd := cmd.ChildSubcmd
	if keySubCmd.Child.Name == "show" {
		return h.processShowCommand(keySubCmd.Child)
	}
	if keySubCmd.Child.Name == "add" {

	}
	return ""
}

func (h *cmdhandler) processShowCommand(show *parser.Subcmd) (result string) {
	if len(show.OptionalArgs) == 0 {
		// Show all keys => values from group
		result = "not implemented yet"
	} else {
		// Show key => value
		result = locales.GetCurrentLocalesStorage().GetDefaultKey(
			show.GetRequiredArg(0).GetValue(),
			show.GetOptionalArg(0).GetValue(),
		)
	}
	return
}

func (h *cmdhandler) processAddCommand(add *parser.Subcmd) (result string) {
	return "not implemented yet"
}
