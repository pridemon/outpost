package models

import (
	"gorm.io/gorm"
)

type RefreshInfo struct {
	Hash         string `gorm:"type:char(32);not null;uniqueIndex:refresh_info_idx;"`
	RefreshToken string `gorm:"type:char(18);not null;uniqueIndex:refresh_info_idx;"`

	gorm.Model
}
