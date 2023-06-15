package locales

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"
	"os"
	"strings"
)

type LocalesStorage struct {
	locales map[string]Locale
}

func (l *LocalesStorage) findLocalesInFolder(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	l.locales = make(map[string]Locale)
	for _, file := range files {
		if file.IsDir() && !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		var locale = Locale{}
		localeLoadErr := locale.Load(path + "/" + file.Name())
		if localeLoadErr != nil {
			logger.Logger().Warn("Can't load locale file %s: %s", file.Name(), localeLoadErr.Error())
			continue
		}
		l.locales[locale.Name] = locale
	}
	l.locales["default"] = defaultLocale

	return nil
}

func (l *LocalesStorage) GetDefaultKey(group, key string) string {
	defaultLocaleFromConfig := config.CurrentConfig.DefaultLocale
	if defaultLocaleFromConfig == "" {
		defaultLocaleFromConfig = "default"
	}
	defaultLocale := l.locales["default"]
	var locale Locale
	var ok bool

	if locale, ok = l.locales[defaultLocaleFromConfig]; !ok {
		value, _ := defaultLocale.GetByKey("general", "locale_not_found")
		return value
	}

	value, err := locale.GetByKey(group, key)
	if err != nil {
		errorKey := ""
		if err == ErrNoGroupFound {
			errorKey = "group_not_found"
		} else {
			errorKey = "key_not_found"
		}
		value, _ := defaultLocale.GetByKey("general", errorKey)
		return value
	}

	return value
}

func (l *LocalesStorage) GetKeyFromLocale(localeName, group, key string) string {
	var locale Locale
	var ok bool

	if locale, ok = l.locales[localeName]; !ok {
		return l.GetDefaultKey(group, key)
	}

	value, err := locale.GetByKey(group, key)
	if err != nil {
		if err == ErrNoGroupFound {
			return l.GetDefaultKey("general", "group_not_found")
		} else {
			return l.GetDefaultKey("general", "key_not_found")
		}
	}

	return value
}
