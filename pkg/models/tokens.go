package models

import (
	"gorm.io/gorm"
)

type Token struct {
	Hash         string `gorm:"type:char(32);not null;uniqueIndex:token_idx;"`
	RefreshToken string `gorm:"type:char(18);not null;uniqueIndex:token_idx;"`

	gorm.Model
}
