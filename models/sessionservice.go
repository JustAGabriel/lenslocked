package models

import (
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
	s := &Session{
		UserID:    userId,
		TokenHash: util.GetSessionToken(),
	}
	_ = ss.db.Model(&Session{}).Create(s)

	return s, nil
}
