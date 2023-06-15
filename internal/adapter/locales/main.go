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

var currentLocalesStorage *LocalesStorage

func New() {
	currentLocalesStorage = &LocalesStorage{}
	err := currentLocalesStorage.findLocalesInFolder(config.CurrentConfig.LocalesFolder)
	if err != nil {
		logger.Logger().Error("Can't process locales folder! Only default locale loaded!")
	}
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
