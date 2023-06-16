package interfaces

import "adblock_bot/internal/repository/models"

type RepositoryManager interface {
	GetVerifierRuleRepository() VerifierRuleRepository
}

type VerifierRuleRepository interface {
	GetAll() ([]models.VerifierRule, error)
	GetByAuthor(author string) ([]models.VerifierRule, error)
	Create(rule models.VerifierRule) error
	RemoveByName(name string) error
}
