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

func (prs *PasswordResetService) DeletePasswordReset(pwr PasswordReset) {
	pwResetDeletionResult := prs.db.Unscoped().Model(&PasswordReset{}).Delete(&pwr)
	if pwResetDeletionResult.Error != nil {
		log.Default().Printf("error while deleting existing pw reset token: %s", pwResetDeletionResult.Error)
	}
}

func (prs *PasswordResetService) GetPasswordReset(usr User) (PasswordReset, error) {
	var existingPwReset PasswordReset
	dbQueryResult := prs.db.Where(&PasswordReset{UserID: usr.ID}).First(&existingPwReset)
	if dbQueryResult.Error != nil {
		return PasswordReset{}, nil
	}

	if existingPwReset.HasExpired() {
		prs.DeletePasswordReset(existingPwReset)
		return PasswordReset{}, errors.New("password reset has expired")
	}

	return existingPwReset, nil
}

func (prs *PasswordResetService) CreatePasswordReset(email string) (PasswordReset, error) {
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

	existingPwReset, err := prs.GetPasswordReset(usr)
	if err == nil {
		prs.DeletePasswordReset(existingPwReset)
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
