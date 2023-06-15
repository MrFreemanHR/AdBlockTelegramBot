package repository

import (
	"adblock_bot/infrastructure/mysql"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/repository/models"
	"errors"

	"gorm.io/gorm"
)

type verifierRuleRepository struct {
	db mysql.MySql
}

var (
	ErrEmptyName = errors.New("empty name field")
	ErrNotFound  = errors.New("rule with this name not found")
)

func NewVerifierRuleRepository(db mysql.MySql) interfaces.VerifierRuleRepository {
	return &verifierRuleRepository{db: db}
}

func (r *verifierRuleRepository) GetAll() ([]models.VerifierRule, error) {
	var rules []models.VerifierRule

	err := r.db.DB.Find(&rules).Error
	if err != nil {
		return nil, err
	}

	return rules, nil
}

func (r *verifierRuleRepository) GetByAuthor(author string) ([]models.VerifierRule, error) {
	var rules []models.VerifierRule

	err := r.db.DB.Where(models.VerifierRule{Author: author}).Find(&rules).Error
	if err != nil {
		return nil, err
	}

	return rules, nil
}

func (r *verifierRuleRepository) Create(rule models.VerifierRule) error {
	if rule.Name == "" {
		return ErrEmptyName
	}
	return r.db.DB.Create(&rule).Error
}

func (r *verifierRuleRepository) RemoveByName(name string) error {
	var rule models.VerifierRule

	err := r.db.DB.Where(models.VerifierRule{Name: name}).First(&rule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else {
			return err
		}
	}

	return r.db.DB.Delete(&rule).Error
}
