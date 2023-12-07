package repository

import (
	"time"

	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/models"
	"gorm.io/gorm"
)

type TokensRepository struct {
	di.Inject

	DB *gorm.DB
}

func (repo *TokensRepository) FindToken(hash string) (*models.Token, error) {
	var token *models.Token

	err := repo.DB.
		Where(models.Token{Hash: hash}).
		First(&token).
		Error

	return token, err
}

func (repo *TokensRepository) Insert(token *models.Token) error {
	return repo.DB.Create(&token).Error
}

func (repo *TokensRepository) DeleteByHash(hash string) error {
	return repo.DB.
		Unscoped().
		Where("hash = ?", hash).
		Delete(&models.Token{}).
		Error
}

func (repo *TokensRepository) DeleteWithCreatedAtLT(checkTime time.Time) (int64, error) {
	res := repo.DB.
		Unscoped().
		Where("created_at < ?", checkTime).
		Delete(&models.Token{})

	return res.RowsAffected, res.Error
}
