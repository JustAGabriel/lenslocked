package models

import "gorm.io/gorm"

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) UserService {
	db.AutoMigrate(&Session{})
	return UserService{
		db: db,
	}
}
