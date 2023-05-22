package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	UserID    uint   `gorm:"primaryKey;unique;not null"`
	TokenHash string `gorm:"unique;not null"`
}
