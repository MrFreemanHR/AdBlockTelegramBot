package locales

var defaultLocale = Locale{
	Name: "default",
	Keys: map[string]map[string]string{
		"general": {
			"cmd_not_found":    "Command not found!",
			"cmd_not_correct":  "Command corrupted!",
			"cmd_arg_error":    "Can't process necessary arguments for command!",
			"key_not_found":    "Key not found!",
			"locale_not_found": "Locale not found!",
			"group_not_found":  "Group in locale not found!",
		},
		"locales": {
			"locales_help": `	Contol all locales in your bot. 
Unfortunately there are no other locales now, so this command will not work correctly :(
			`,
			"locales_group_show":  "Group name: %s\n%s",
			"locales_key_added":   "Add key \"%s\" to group \"%s\" with value \"%s\"",
			"locale_saving_error": "Error while saving current locale: %s",
			"locale_saved":        "Locale saved successfully!",
		},
		"verifier": {
			"text_regexp_invalid":       "Text regexp is invalid: %s",
			"rule_created":              "New rule with name \"%s\" successfully created!",
			"rule_creating_error":       "Can't create rule!",
			"rule_removed":              "Rule with name \"%s\" removed successfully!",
			"rule_removing_error":       "Can't remove this rule!",
			"rule_get_error":            "Can't get rules right now due to internal error",
			"rules_with_author_head":    "=== <i>Rules for %s message's author:</i> ===",
			"rules_without_author_head": "=== <i>All rules:</i> ===",
			"rule_head":                 "<b>Rule</b> %s:",
			"rule_raw_created_at":       "<b>Created at</b>: %s",
			"rule_raw_author":           "<b>Author</b>: %s",
			"rule_raw_audio_note":       "<b>Audio note</b>: %s",
			"rule_raw_video_note":       "<b>Video note</b>: %s",
			"rule_raw_photo":            "<b>Photo</b>: %s",
			"rule_raw_text":             "<b>Text regexp</b>: %s",
			"rule_raw_locale_group":     "<b>Group from locale to reply</b>: %s",
			"rule_raw_key_group":        "<b>Key from group to reply</b>: %s",
			"rule_random_key":           "Random key from this group",
		},
	},
}
