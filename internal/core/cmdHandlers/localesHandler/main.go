package localeshandler

import (
	"adblock_bot/internal/adapter/locales"
	"adblock_bot/internal/adapter/parser"
	"adblock_bot/internal/core/interfaces"
	"fmt"
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
		return h.processAddCommand(keySubCmd.Child)
	}
	return ""
}

func (h *cmdhandler) processShowCommand(show *parser.Subcmd) (result string) {
	if len(show.OptionalArgs) == 0 {
		// Show all keys => values from group
		keys, err := locales.GetCurrentLocalesStorage().GetAllKeysFromDefaultLocale(
			show.GetRequiredArg(0).GetValue(),
		)
		if err == locales.ErrGroupNotFound {
			result = locales.GetCurrentLocalesStorage().GetDefaultKey("general", "group_not_found")
		} else if err == locales.ErrNoGroupFound {
			result = locales.GetCurrentLocalesStorage().GetDefaultKey("general", "locale_not_found")
		} else {
			var bufferStr string
			for key, value := range keys {
				bufferStr += fmt.Sprintf("%s => %s\n", key, value)
			}
			result = fmt.Sprintf(
				locales.GetCurrentLocalesStorage().GetDefaultKey("locales", "locales_group_show"),
				show.GetRequiredArg(0).GetValue(),
				bufferStr,
			)
		}
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
	locales.GetCurrentLocalesStorage().AddKeyToDefaultLocale(
		add.GetRequiredArg(0).GetValue(),
		add.GetRequiredArg(1).GetValue(),
		add.GetRequiredArg(2).GetValue(),
	)
	return fmt.Sprintf(
		locales.GetCurrentLocalesStorage().GetDefaultKey("locales", "locales_key_added"),
		add.GetRequiredArg(1).GetValue(),
		add.GetRequiredArg(0).GetValue(),
		add.GetRequiredArg(2).GetValue(),
	)
}
