package verifier

import (
	"adblock_bot/infrastructure/mysql"
	"adblock_bot/internal/adapter/locales"
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/core/entity"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/repository"
	"adblock_bot/internal/repository/models"
	"adblock_bot/internal/transport"
	"regexp"
)

type handler struct {
	api         *transport.TelegramAPI
	verRuleRepo interfaces.VerifierRuleRepository
}

func New(api *transport.TelegramAPI, db *mysql.MySql) interfaces.MessageHandler {
	repoManager := repository.NewRepositoryManager(*db)
	verRuleRepo := repoManager.GetVerifierRuleRepository()

	return &handler{
		api:         api,
		verRuleRepo: verRuleRepo,
	}
}

func (h *handler) ProcessMessage(event *entity.TelegramMessage) bool {
	if event.Chat.IsPrivate() {
		return false
	}
	user := event.From
	if user == nil {
		return false
	}
	rules, err := h.verRuleRepo.GetByAuthor(user.UserName)
	if err != nil {
		logger.Logger().Warn("Can't get rules from db for user %s: %s", user.UserName, err.Error())
		return false
	}
	reply := entity.NewMessage(event.Chat.ID, "")
	reply.ReplyToMessage = event
	for _, rule := range rules {
		if h.processRule(event, rule) {
			// Process keys from local
			if rule.LocaleGroup == "" {
				h.api.RemoveMessage(event.Chat.ID, event.MessageID)
				return true
			}
			if rule.LocaleKey != "" {
				reply.Text = locales.GetCurrentLocalesStorage().GetDefaultKey(rule.LocaleGroup, rule.LocaleKey)
			} else {
				reply.Text = locales.GetCurrentLocalesStorage().GetRandomKeyFromCurrentLocale(rule.LocaleGroup)
			}
			h.api.SendMessage(reply)
			h.api.RemoveMessage(event.Chat.ID, event.MessageID)
			return true
		}
	}
	return false
}

func (h *handler) processRule(event *entity.TelegramMessage, rule models.VerifierRule) bool {
	// Check for deny audio note
	if event.Audio != nil && !rule.AudioNote {
		return true
	}
	// Check for deny video note
	if event.VideoNote != nil && !rule.VideoNote {
		return true
	}
	// Check for deny photo
	if event.Photo != nil && !rule.Photo {
		return true
	}
	// Check text by regexp
	if event.Text != "" && rule.Text != "" {
		r, err := regexp.Compile(rule.Text)
		if err != nil {
			logger.Logger().Warn("Invalid regexp in rule with id %d: %s", rule.Id, err.Error())
			return false
		}
		return r.MatchString(event.Text)
	}
	return false
}
