package verifier

import (
	"adblock_bot/infrastructure/sqlite"
	"adblock_bot/internal/adapter/locales"
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/repository"
	"adblock_bot/internal/repository/models"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler struct {
	bot         *tgbotapi.BotAPI
	verRuleRepo interfaces.VerifierRuleRepository
}

func New(bot *tgbotapi.BotAPI, db *sqlite.SQLite) interfaces.MessageHandler {
	repoManager := repository.NewRepositoryManager(*db)
	verRuleRepo := repoManager.GetVerifierRuleRepository()

	return &handler{
		bot:         bot,
		verRuleRepo: verRuleRepo,
	}
}

func (h *handler) ProcessMessage(event *tgbotapi.Update) bool {
	user := event.Message.From
	if user == nil {
		return false
	}
	rules, err := h.verRuleRepo.GetByAuthor(user.UserName)
	if err != nil {
		logger.Logger().Warn("Can't get rules from db for user %s: %s", user.UserName, err.Error())
		return false
	}
	reply := tgbotapi.NewMessage(event.Message.Chat.ID, "")
	reply.ReplyToMessageID = event.Message.MessageID
	del := tgbotapi.NewDeleteMessage(event.Message.Chat.ID, event.Message.MessageID)
	for _, rule := range rules {
		if h.processRule(event, rule) {
			// Process keys from local
			if rule.LocaleGroup == "" {
				h.bot.Send(del)
				return true
			}
			if rule.LocaleKey != "" {
				reply.Text = locales.GetCurrentLocalesStorage().GetDefaultKey(rule.LocaleGroup, rule.LocaleKey)
			} else {
				reply.Text = locales.GetCurrentLocalesStorage().GetRandomKeyFromCurrentLocale(rule.LocaleGroup)
			}
			h.bot.Send(reply)
			h.bot.Send(del)
			return true
		}
	}
	return false
}

func (h *handler) processRule(event *tgbotapi.Update, rule models.VerifierRule) bool {
	// Check for deny audio note
	if event.Message.Audio != nil && !rule.AudioNote {
		return true
	}
	// Check for deny video note
	if event.Message.VideoNote != nil && !rule.VideoNote {
		return true
	}
	// Check for deny photo
	if event.Message.Photo != nil && !rule.Photo {
		return true
	}
	// Check text by regexp
	if event.Message.Text != "" && rule.Text != "" {
		r, err := regexp.Compile(rule.Text)
		if err != nil {
			logger.Logger().Warn("Invalid regexp in rule with id %d: %s", rule.Id, err.Error())
			return false
		}
		return r.MatchString(event.Message.Text)
	}
	return false
}
