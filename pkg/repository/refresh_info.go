package repository

import (
	"time"

	"github.com/goava/di"
	"github.com/pridemon/outpost/pkg/models"
	"gorm.io/gorm"
)

type RefreshInfoRepository struct {
	di.Inject

	DB *gorm.DB
}

func (repo *RefreshInfoRepository) Find(hash string) (*models.RefreshInfo, error) {
	var token *models.RefreshInfo

	err := repo.DB.
		Where(models.RefreshInfo{Hash: hash}).
		First(&token).
		Error

	return token, err
}

func (repo *RefreshInfoRepository) Insert(token *models.RefreshInfo) error {
	return repo.DB.Create(&token).Error
}

func (repo *RefreshInfoRepository) DeleteByHash(hash string) error {
	return repo.DB.
		Unscoped().
		Where("hash = ?", hash).
		Delete(&models.RefreshInfo{}).
		Error
}

func (repo *RefreshInfoRepository) DeleteWithCreatedAtLT(checkTime time.Time) (int64, error) {
	res := repo.DB.
		Unscoped().
		Where("created_at < ?", checkTime).
		Delete(&models.RefreshInfo{})

	return res.RowsAffected, res.Error
}
