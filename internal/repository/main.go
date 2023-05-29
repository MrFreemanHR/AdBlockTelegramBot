package repository

import (
	"adblock_bot/infrastructure/sqlite"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/repository/models"
)

type repositoryManager struct {
	db sqlite.SQLite
}

func NewRepositoryManager(db sqlite.SQLite) interfaces.RepositoryManager {
	var repoManager = &repositoryManager{
		db: db,
	}
	repoManager.db.DB.AutoMigrate(&models.VerifierRule{})
	return repoManager
}

func (m *repositoryManager) GetVerifierRuleRepository() interfaces.VerifierRuleRepository {
	return NewVerifierRuleRepository(m.db)
}
