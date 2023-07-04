package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type PasswordResetService struct {
	db            *gorm.DB
	bytesPerToken int
	duration      time.Duration
}

func NewPasswordResetService(db *gorm.DB, bytesPerToken int, duration time.Duration) *PasswordResetService {
	return &PasswordResetService{
		db:            db,
		bytesPerToken: bytesPerToken,
		duration:      duration,
	}
}

func (prs *PasswordResetService) GetPasswordReset(email string) (PasswordReset, error) {
	return PasswordReset{}, errors.New("not implemented")
}

func (prs *PasswordResetService) GetUserByToken(token string) (User, error) {
	return User{}, errors.New("not implemented")
}
