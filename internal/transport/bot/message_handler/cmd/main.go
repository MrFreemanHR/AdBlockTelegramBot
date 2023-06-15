package cmd

import (
	"adblock_bot/infrastructure/sqlite"
	"adblock_bot/internal/adapter/locales"
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/adapter/parser"
	localeshandler "adblock_bot/internal/core/cmdHandlers/localesHandler"
	verifierhandler "adblock_bot/internal/core/cmdHandlers/verifierHandler"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/repository"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler struct {
	bot    *tgbotapi.BotAPI
	lch    interfaces.CmdHandler
	vfh    interfaces.CmdHandler
	parser *parser.Parser
}

func New(bot *tgbotapi.BotAPI, db *sqlite.SQLite) *handler {
	repoManager := repository.NewRepositoryManager(*db)
	verRuleRepo := repoManager.GetVerifierRuleRepository()

	return &handler{
		bot:    bot,
		lch:    localeshandler.New(),
		vfh:    verifierhandler.New(verRuleRepo),
		parser: parser.New(),
	}
}

func (h *handler) ProcessMessage(event *tgbotapi.Update) bool {
	messageText := event.Message.Text

	if len(messageText) > 0 && messageText[0] == '/' {
		if event.Message.Chat.IsSuperGroup() && !strings.Contains(messageText, h.bot.Self.UserName) {
			return false
		}

		if strings.Contains(messageText, "@") {
			messageText = messageText[:strings.Index(messageText, "@")]
		}
		cmd, err := h.parser.Parse(messageText)
		reply := tgbotapi.NewMessage(event.Message.Chat.ID, "")
		var replyText string
		if err != nil {
			replyText = h.processErrorFromParser(err)
		} else {
			if cmd.Name == "locales" {
				replyText = h.lch.ProcessCommand(cmd)
			}
			if cmd.Name == "verifier" {
				replyText = h.vfh.ProcessCommand(cmd)
				reply.ParseMode = "html"
			}
		}
		reply.Text = replyText
		reply.ReplyToMessageID = event.Message.MessageID
		_, err = h.bot.Send(reply)
		if err != nil {
			logger.Logger().Error("Error while sending reply message: %s", err.Error())
		}
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
