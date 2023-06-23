package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	gorm.Model
	UserID    uint
	Token     string
	ExpiresAt time.Time
}
