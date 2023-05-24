package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	UserID uint   `gorm:"primaryKey;unique;not null"`
	Token  string `gorm:"unique;not null"`
}
