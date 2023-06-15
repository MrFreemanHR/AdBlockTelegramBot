package verifierhandler

import (
	"adblock_bot/internal/adapter/locales"
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/adapter/parser"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/repository/models"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type cmdhandler struct {
	verifierRuleRepository interfaces.VerifierRuleRepository
}

func New(verifierRuleRepository interfaces.VerifierRuleRepository) interfaces.CmdHandler {
	return &cmdhandler{
		verifierRuleRepository: verifierRuleRepository,
	}
}

func (h *cmdhandler) ProcessCommand(cmd *parser.Cmd) string {
	subcmd := cmd.ChildSubcmd
	if subcmd.Name == "add" {
		return h.processAddCommand(subcmd)
	}
	if subcmd.Name == "remove" {
		return h.processRemoveCommand(subcmd)
	}
	if subcmd.Name == "list" {
		return h.processListCommand(subcmd)
	}
	return ""
}

func (h *cmdhandler) processArg(rawArg string) (key, value string) {
	delimiterIndex := strings.Index(rawArg, ":")
	if delimiterIndex == -1 {
		return
	}
	key = strings.TrimSpace(rawArg[:delimiterIndex])
	value = strings.TrimSpace(rawArg[delimiterIndex+1:])
	return
}

func (h *cmdhandler) convertToBool(arg string) bool {
	lowerArg := strings.ToLower(arg)
	return lowerArg == "allow"
	// If not allowed clearly - deny this
}

func (h *cmdhandler) convertToString(arg bool) string {
	if arg {
		return "allow"
	} else {
		return "deny"
	}
}

func (h *cmdhandler) processAddCommand(cmd *parser.Subcmd) string {
	ruleName := cmd.GetRequiredArg(0).GetValue()
	rawAuthor := cmd.GetRequiredArg(1).GetValue()
	authorKey, authorValue := h.processArg(rawAuthor)
	if authorKey == "" || authorValue == "" || authorKey != "author" {
		return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_arg_error")
	}

	var rule = models.VerifierRule{
		Name:      ruleName,
		Author:    authorValue,
		AudioNote: true,
		VideoNote: true,
		Photo:     true,
	}
	for _, optionalArg := range cmd.OptionalArgs {
		rawArg := optionalArg.GetValue()
		argKey, argValue := h.processArg(rawArg)
		if argKey == "" || argValue == "" {
			continue
		}
		if argKey == "audio_note" {
			rule.AudioNote = h.convertToBool(argValue)
		}
		if argKey == "video_note" {
			rule.VideoNote = h.convertToBool(argValue)
		}
		if argKey == "photo" {
			rule.Photo = h.convertToBool(argValue)
		}
		if argKey == "message" {
			rawValues := strings.Split(argValue, ",")
			if len(rawValues) == 0 {
				return locales.GetCurrentLocalesStorage().GetDefaultKey("general", "cmd_arg_error")
			}
			if len(rawValues) == 1 {
				rule.LocaleGroup = rawValues[0]
			} else {
				rule.LocaleKey = rawValues[1]
			}
		}
		if argKey == "text" {
			_, err := regexp.Compile(argValue)
			if err != nil {
				return fmt.Sprintf(
					locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "text_regexp_invalid"),
					err.Error(),
				)
			}
			rule.Text = argValue
		}
	}
	err := h.verifierRuleRepository.Create(rule)
	if err != nil {
		logger.Logger().Error("Can't create rule: %s", err.Error())
		return locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_creating_error")
	}
	return fmt.Sprintf(
		locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_created"),
		rule.Name,
	)
}

func (h *cmdhandler) processRemoveCommand(cmd *parser.Subcmd) string {
	name := cmd.GetRequiredArg(0).GetValue()
	err := h.verifierRuleRepository.RemoveByName(name)
	if err != nil {
		logger.Logger().Error("Can't remove rule: %s", err.Error())
		return locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_removing_error")
	}
	return fmt.Sprintf(
		locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_removed"),
		name,
	)
}

func (h *cmdhandler) processListCommand(cmd *parser.Subcmd) string {
	var jsonFlag bool
	var rules []models.VerifierRule
	var err error

	for _, optionalArg := range cmd.OptionalArgs {
		rawArg := optionalArg.GetValue()
		if strings.ToLower(rawArg) == "json" {
			jsonFlag = true
		}
		argKey, argValue := h.processArg(rawArg)
		if argKey == "author" {
			rules, err = h.verifierRuleRepository.GetByAuthor(argValue)
			if err != nil {
				logger.Logger().Error("Can't get rules: %s", err.Error())
				return locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_get_error")
			}
			return h.printableRules(rules, jsonFlag, argValue)
		}
	}
	rules, err = h.verifierRuleRepository.GetAll()
	if err != nil {
		logger.Logger().Error("Can't get rules: %s", err.Error())
		return locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_get_error")
	}
	return h.printableRules(rules, jsonFlag)
}

func (h *cmdhandler) printableRules(rules []models.VerifierRule, jsonFlag bool, author ...string) (result string) {
	if len(author) == 0 {
		result = locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rules_without_author_head")
	} else {
		result = fmt.Sprintf(
			locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rules_with_author_head"),
			author[0],
		)
	}
	result += "\n"
	for _, rule := range rules {
		result += (fmt.Sprintf(
			locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_head")+"\n",
			rule.Name,
		))

		if jsonFlag {
			data, _ := json.MarshalIndent(
				rule,
				" ",
				"	",
			)
			result += fmt.Sprintf("<pre>%s</pre>", string(data)) + "\n"
		} else {
			result += (fmt.Sprintf(
				locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_raw_created_at"),
				rule.CreatedAt.Format(time.RFC822),
			) + "\n")
			result += (fmt.Sprintf(
				locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_raw_author"),
				rule.Author,
			) + "\n")
			result += (fmt.Sprintf(
				locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_raw_audio_note"),
				h.convertToString(rule.AudioNote),
			) + "\n")
			result += (fmt.Sprintf(
				locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_raw_video_note"),
				h.convertToString(rule.VideoNote),
			) + "\n")
			result += (fmt.Sprintf(
				locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_raw_photo"),
				h.convertToString(rule.Photo),
			) + "\n")
			result += (fmt.Sprintf(
				locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_raw_text"),
				rule.Text,
			) + "\n")
			if rule.LocaleGroup != "" {
				result += (fmt.Sprintf(
					locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_raw_locale_group"),
					rule.LocaleGroup,
				) + "\n")
				if rule.LocaleKey == "" {
					result += (fmt.Sprintf(
						locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_raw_key_group"),
						locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_random_key"),
					) + "\n")
				} else {
					result += (fmt.Sprintf(
						locales.GetCurrentLocalesStorage().GetDefaultKey("verifier", "rule_raw_key_group"),
						rule.LocaleKey,
					) + "\n")
				}
			}
		}
		result += "=== === === === ===\n"
	}
	return result
}
