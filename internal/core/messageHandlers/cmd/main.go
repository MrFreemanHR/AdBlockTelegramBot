package cmd

import (
	"adblock_bot/infrastructure/mysql"
	"adblock_bot/internal/adapter/locales"
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/adapter/parser"
	localeshandler "adblock_bot/internal/core/cmdHandlers/localesHandler"
	verifierhandler "adblock_bot/internal/core/cmdHandlers/verifierHandler"
	"adblock_bot/internal/core/entity"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/repository"
	"adblock_bot/internal/transport"
	"strings"
)

type handler struct {
	api    *transport.TelegramAPI
	lch    interfaces.CmdHandler
	vfh    interfaces.CmdHandler
	parser *parser.Parser
}

func New(api *transport.TelegramAPI, db *mysql.MySql) interfaces.MessageHandler {
	repoManager := repository.NewRepositoryManager(*db)
	verRuleRepo := repoManager.GetVerifierRuleRepository()

	return &handler{
		api:    api,
		lch:    localeshandler.New(),
		vfh:    verifierhandler.New(verRuleRepo),
		parser: parser.New(),
	}
}

func (h *handler) ProcessMessage(event *entity.TelegramMessage) bool {
	messageText := event.Text

	if len(messageText) > 0 && messageText[0] == '/' {
		if event.Chat.IsSuperGroup() && !strings.Contains(messageText, h.api.Self().UserName) {
			return false
		}

		if strings.Contains(messageText, "@") {
			messageText = messageText[:strings.Index(messageText, "@")]
		}
		cmd, err := h.parser.Parse(messageText)
		reply := entity.NewMessage(event.Chat.ID, "")
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
		reply.ReplyToMessage = event
		err = h.api.SendMessage(reply)
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
