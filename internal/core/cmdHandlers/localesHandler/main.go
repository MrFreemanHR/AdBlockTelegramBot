package localeshandler

import (
	"adblock_bot/internal/adapter/locales"
	"adblock_bot/internal/adapter/logger"
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
	SubCmd := cmd.ChildSubcmd
	if SubCmd.Name == "key" {
		if SubCmd.Child.Name == "show" {
			return h.processShowCommand(SubCmd.Child)
		}
		if SubCmd.Child.Name == "add" {
			return h.processAddCommand(SubCmd.Child)
		}
	}
	if SubCmd.Name == "locale" {
		return h.processLocaleCommand()
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

func (h *cmdhandler) processLocaleCommand() (result string) {
	var err error
	var locale locales.Locale
	locale, err = locales.GetCurrentLocalesStorage().GetCurrentLocale()
	if err != nil {
		locale = locales.GetCurrentLocalesStorage().GetDefaultLocale()
		logger.Logger().Warn("Can't find current locale from config")
	}
	err = locales.GetCurrentLocalesStorage().SaveLocale(locale)
	if err != nil {
		return fmt.Sprintf(locales.GetCurrentLocalesStorage().GetDefaultKey("locales", "locale_saving_error"), err.Error())
	}
	return locales.GetCurrentLocalesStorage().GetDefaultKey("locales", "locale_saved")
}
