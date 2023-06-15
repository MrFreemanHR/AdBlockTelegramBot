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
	},
}
