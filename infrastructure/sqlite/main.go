package sqlite

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLite struct {
	DB *gorm.DB
}

func New() *SQLite {
	db, err := gorm.Open(
		sqlite.Open(config.CurrentConfig.SQLiteDSN),
	)
	if err != nil {
		logger.Logger().Fatal("Can't connect to internal db: %s", err.Error())
		return nil
	}
	return &SQLite{
		DB: db,
	}
}
