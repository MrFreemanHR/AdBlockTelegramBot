package cmd

import (
	"adblock_bot/internal/adapter/locales"
	"adblock_bot/internal/adapter/parser"
	localeshandler "adblock_bot/internal/core/cmdHandlers/localesHandler"
	"adblock_bot/internal/core/interfaces"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler struct {
	bot    *tgbotapi.BotAPI
	lch    interfaces.CmdHandler
	parser *parser.Parser
}

func New(bot *tgbotapi.BotAPI) *handler {
	return &handler{
		bot:    bot,
		lch:    localeshandler.New(),
		parser: parser.New(),
	}
}

func (h *handler) ProcessMessage(event *tgbotapi.Update) bool {
	messageText := event.Message.Text

	if strings.Contains(messageText, "/locales") {
		cmd, err := h.parser.Parse(messageText)
		var replyText string
		if err != nil {
			replyText = h.processErrorFromParser(err)
		} else {
			replyText = h.lch.ProcessCommand(cmd)
		}
		reply := tgbotapi.NewMessage(event.Message.Chat.ID, replyText)
		reply.ReplyToMessageID = event.Message.MessageID
		h.bot.Send(reply)
		return true
	}

	return false
}

func (h *handler) processErrorFromParser(err error) string {
	if err == parser.ErrNotKnownCommand || err == parser.ErrNoCommand {
		return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_not_found")
	}
	if err == parser.ErrParsingSubCmd || err == parser.ErrNotKnownSubCommand {
		return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_not_correct")
	}
	if err == parser.ErrParsingRequiredArg || err == parser.ErrParsingOptionalArg {
		return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_arg_error")
	}
	return err.Error()
}
