package mysql

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySql struct {
	DB *gorm.DB
}

func New() *MySql {
	db, err := gorm.Open(
		mysql.Open(config.CurrentConfig.MySQlDSN),
	)
	if err != nil {
		logger.Logger().Fatal("Can't connect to internal db: %s", err.Error())
		return nil
	}
	return &MySql{
		DB: db,
	}
}
