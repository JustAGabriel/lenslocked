package models

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/justagabriel/lenslocked/util"
	"gorm.io/gorm"
)

const (
	SessionTokenLength = 32
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) SessionService {
	db.AutoMigrate(&Session{})
	return SessionService{
		db: db,
	}
}

func (ss *SessionService) GetNewSession(userId uint) (*Session, error) {
	token := util.GetSessionToken()
	hashedTokenBytes := sha256.Sum256([]byte(token))
	hashedToken := base64.URLEncoding.EncodeToString(hashedTokenBytes[:])
	s := &Session{
		UserID: userId,
		Token:  hashedToken,
	}

	_ = ss.db.Model(&Session{}).Create(s)

	return s, nil
}
