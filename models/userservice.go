package models

import (
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GeneratePasswordHash(pw string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	db.AutoMigrate(&User{})
	return UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(u User) {
	u.Email = strings.ToLower(u.Email)
	_ = s.db.Model(&User{}).Create(&u)
}
