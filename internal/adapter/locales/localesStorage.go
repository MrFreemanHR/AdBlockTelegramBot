package locales

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"
	"errors"
	"math/rand"
	"os"
	"strings"
	"time"
)

type LocalesStorage struct {
	locales map[string]Locale
}

var (
	ErrLocaleNotFound = errors.New("locale with this name not found")
	ErrGroupNotFound  = errors.New("group with this name not found")
)

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

func (l *LocalesStorage) GetCurrentLocale() (Locale, error) {
	defaultLocaleFromConfig := config.CurrentConfig.DefaultLocale
	if defaultLocaleFromConfig == "" {
		defaultLocaleFromConfig = "default"
	}

	if locale, ok := l.locales[defaultLocaleFromConfig]; ok {
		return locale, nil
	}
	return Locale{}, ErrLocaleNotFound
}

func (l *LocalesStorage) GetDefaultLocale() Locale {
	return l.locales["default"]
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

func (l *LocalesStorage) GetRandomKey(localeName, group string) string {
	defaultLocale := l.locales["default"]
	var locale Locale
	var ok bool

	if locale, ok = l.locales[localeName]; !ok {
		value, _ := defaultLocale.GetByKey("general", "locale_not_found")
		return value
	}

	if locale.Keys[group] == nil || len(locale.Keys[group]) == 0 {
		value, _ := defaultLocale.GetByKey("general", "group_not_found")
		return value
	}

	keys := make([]string, 0, len(locale.Keys[group]))
	for k := range locale.Keys[group] {
		keys = append(keys, k)
	}

	rs := rand.NewSource(time.Now().Unix())
	r := rand.New(rs)
	randomInt := r.Intn(len(keys))
	randomKey := keys[randomInt]
	val, _ := locale.GetByKey(group, randomKey)
	return val
}

func (l *LocalesStorage) GetRandomKeyFromCurrentLocale(group string) string {
	defaultLocaleFromConfig := config.CurrentConfig.DefaultLocale
	if defaultLocaleFromConfig == "" {
		defaultLocaleFromConfig = "default"
	}

	return l.GetRandomKey(defaultLocaleFromConfig, group)
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

func (l *LocalesStorage) GetAllKeysFromLocale(localeName, groupName string) (map[string]string, error) {
	var locale Locale
	var group map[string]string
	var ok bool

	if locale, ok = l.locales[localeName]; !ok {
		return nil, ErrLocaleNotFound
	}

	if group, ok = locale.Keys[groupName]; !ok {
		return nil, ErrGroupNotFound
	}

	return group, nil
}

func (l *LocalesStorage) GetAllKeysFromDefaultLocale(groupName string) (map[string]string, error) {
	defaultLocaleFromConfig := config.CurrentConfig.DefaultLocale
	if defaultLocaleFromConfig == "" {
		defaultLocaleFromConfig = "default"
	}
	return l.GetAllKeysFromLocale(defaultLocaleFromConfig, groupName)
}

func (l *LocalesStorage) AddKeyToDefaultLocale(groupName, key, value string) error {
	defaultLocaleFromConfig := config.CurrentConfig.DefaultLocale
	if defaultLocaleFromConfig == "" {
		defaultLocaleFromConfig = "default"
	}

	var locale Locale
	var ok bool
	if locale, ok = l.locales[defaultLocaleFromConfig]; !ok {
		return ErrLocaleNotFound
	}

	locale.SetByKey(groupName, key, value)

	return nil
}

func (l *LocalesStorage) SaveLocale(locale Locale) error {
	path := config.CurrentConfig.LocalesFolder + "/" + locale.Name + ".json"
	return locale.Save(path)
}
