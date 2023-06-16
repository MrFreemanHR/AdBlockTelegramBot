package locales

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"
)

var currentLocalesStorage *LocalesStorage

func New() {
	currentLocalesStorage = &LocalesStorage{}
	err := currentLocalesStorage.findLocalesInFolder(config.CurrentConfig.LocalesFolder)
	if err != nil {
		logger.Logger().Error("Can't process locales folder! Only default locale loaded!")
	}
}

func GetCurrentLocalesStorage() *LocalesStorage {
	return currentLocalesStorage
}
