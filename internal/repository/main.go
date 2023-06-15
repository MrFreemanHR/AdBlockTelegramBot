package repository

import (
	"adblock_bot/infrastructure/mysql"
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/repository/models"
)

type repositoryManager struct {
	db mysql.MySql
}

func NewRepositoryManager(db mysql.MySql) interfaces.RepositoryManager {
	var repoManager = &repositoryManager{
		db: db,
	}
	migrateErr := repoManager.db.DB.AutoMigrate(&models.VerifierRule{})
	if migrateErr != nil {
		logger.Logger().Fatal("Can't auto migrate db: %s", migrateErr.Error())
		return nil
	}
	return repoManager
}

func (m *repositoryManager) GetVerifierRuleRepository() interfaces.VerifierRuleRepository {
	return NewVerifierRuleRepository(m.db)
}
