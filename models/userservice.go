package models

import (
	"fmt"
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

func (s *UserService) GetUserByEmail(email string) (User, error) {
	normalizedEmail := strings.ToLower(email)
	log.Default().Printf("try retrieving user with email '%s'", normalizedEmail)

	var usr User
	res := s.db.Where(&User{Email: normalizedEmail}).First(&usr)

	if res.Error != nil {
		return usr, res.Error
	}

	if usr.Email == "" {
		return usr, fmt.Errorf("could not find user with email '%s'", normalizedEmail)
	}

	return usr, nil
}
