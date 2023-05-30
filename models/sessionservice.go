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

func hashToken(token string) string {
	hashedTokenBytes := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hashedTokenBytes[:])
}

func NewSessionService(db *gorm.DB) SessionService {
	db.AutoMigrate(&Session{})
	return SessionService{
		db: db,
	}
}

func (ss *SessionService) GetNewSession(userId uint) (*Session, error) {
	token := util.GetSessionToken()
	hashedToken := hashToken(token)
	s := &Session{
		UserID: userId,
		Token:  hashedToken,
	}

	_ = ss.db.Model(&Session{}).Create(s)

	return s, nil
}

func (ss *SessionService) GetSessionByToken(unhashedToken string) (Session, error) {
	hashedToken := hashToken(unhashedToken)
	session := Session{
		Token: hashedToken,
	}

	res := ss.db.Where(&session).First(&session)
	if res.Error != nil {
		return Session{}, res.Error
	}

	return session, nil
}

func (ss *SessionService) DeleteSessionByToken(token string) error {
	result := ss.db.Delete(&Session{Token: token})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
