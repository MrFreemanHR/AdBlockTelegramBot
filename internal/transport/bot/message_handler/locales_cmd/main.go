package localescmd

import (
	"adblock_bot/internal/adapter/locales"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler struct {
	bot *tgbotapi.BotAPI
}

func New(bot *tgbotapi.BotAPI) *handler {
	return &handler{bot: bot}
}

func (h *handler) ProcessMessage(event *tgbotapi.Update) bool {
	messageText := event.Message.Text

	if strings.Contains(messageText, "/locales") {
		cmdString := messageText[strings.Index(messageText, "/locales"):]
		reply := tgbotapi.NewMessage(event.Message.Chat.ID, h.processCmd(cmdString))
		reply.ReplyToMessageID = event.Message.MessageID
		h.bot.Send(reply)
		return true
	}

	return false
}

func (h *handler) processCmd(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	cmdParts := strings.Split(cmd, " ")
	if len(cmdParts) == 1 {
		return locales.GetCurrentLocalesStorage().GetDefaultKey("locales", "locales_help")
	}

	var cmdLevel = uint8(1)
	var cmdHeadPart, keyHeadPart bool
	var key string
	var i int

	for i, key = range cmdParts {
		if key == "" {
			continue
		}
		if key == "/locales" {
			if cmdLevel != 1 {
				return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_not_correct")
			}
			cmdHeadPart = true
		}
		if key == "key" {
			if cmdLevel != 2 {
				return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_not_correct")
			}
			keyHeadPart = true
		}
		cmdLevel++
	}

	if cmdHeadPart && !keyHeadPart {
		return locales.GetCurrentLocalesStorage().GetDefaultKey("locales", "locales_help")
	}

	if keyHeadPart {
		return h.processKeyPart(cmdParts[i-1:])
	}

	return strings.Join(cmdParts, "|")
}

func (h *handler) processKeyPart(cmdParts []string) string {
	// locales key {show,add} group [key, value]

	var cmdLevel = uint8(1)
	var showPart, addPart bool
	var groupName, keyName, value string

	fmt.Printf("Keys: %#+v\n", cmdParts)
	for _, key := range cmdParts {
		fmt.Printf("key: %s; count: %d\n", key, cmdLevel)
		if key == "" {
			continue
		}
		if key == "show" {
			if cmdLevel != 1 {
				return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_not_correct")
			}
			showPart = true
		}
		if key == "add" {
			if cmdLevel != 1 {
				return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_not_correct")
			}
			addPart = true
		}
		if cmdLevel == 2 {
			groupName = key
		}
		if cmdLevel == 3 {
			keyName = key
		}
		if cmdLevel == 4 {
			value = key
		}
		cmdLevel++
	}

	if showPart {
		return "show all group here: " + groupName
	}

	if addPart {
		return fmt.Sprintf("add key %s with value %s to group %s", keyName, value, groupName)
	}

	return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_not_correct")
}
