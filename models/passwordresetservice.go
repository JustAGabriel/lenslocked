package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

const (
	tokenBytes                     = 256
	defaultTokenExpirationDuration = time.Minute * 15
)

type PasswordResetService struct {
	userService *UserService
	db          *gorm.DB
}

func NewPasswordResetService(userService *UserService, db *gorm.DB) *PasswordResetService {
	db.AutoMigrate(&PasswordReset{})
	return &PasswordResetService{
		userService: userService,
		db:          db,
	}
}

func (prs *PasswordResetService) GetPasswordReset(email string) (PasswordReset, error) {
	usr, err := prs.userService.GetUserByEmail(email)
	if err != nil {
		return PasswordReset{}, err
	}

	tokenBytes := make([]byte, tokenBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return PasswordReset{}, err
	}
	token := base64.URLEncoding.Strict().EncodeToString(tokenBytes)

	pwReset := PasswordReset{
		UserID:    usr.ID,
		Token:     token,
		ExpiresAt: time.Now().UTC().Add(defaultTokenExpirationDuration),
	}

	var existingPwReset PasswordReset
	_ = prs.db.Where(&PasswordReset{UserID: usr.ID}).First(&existingPwReset)
	if existingPwReset.Token != "" {
		pwResetDeletionResult := prs.db.Unscoped().Model(&PasswordReset{}).Where(&PasswordReset{UserID: usr.ID}).Delete(&PasswordReset{})
		if pwResetDeletionResult.Error != nil {
			log.Default().Printf("error while deleting existing pw reset token: %s", pwResetDeletionResult.Error)
		}
	}

	pwResetCreationResult := prs.db.Model(&PasswordReset{}).Create(&pwReset)
	if pwResetCreationResult.Error != nil {
		return PasswordReset{}, pwResetCreationResult.Error
	}

	return pwReset, nil
}

func (prs *PasswordResetService) GetUserByToken(token string) (User, error) {
	return User{}, errors.New("not implemented")
}
