package models

import (
	"fmt"
	"log"
	"strings"

	m "github.com/justagabriel/lenslocked/models"
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
	db.AutoMigrate(&m.User{})
	return UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(u m.User) {
	u.Email = strings.ToLower(u.Email)
	_ = s.db.Model(&m.User{}).Create(&u)
}

func (s *UserService) GetUserByEmail(email string) (m.User, error) {
	normalizedEmail := strings.ToLower(email)
	log.Default().Printf("try retrieving user with email '%s'", normalizedEmail)

	var usr m.User
	res := s.db.Where(&m.User{Email: normalizedEmail}).First(&usr)
	if res.Error != nil {
		return usr, res.Error
	}

	if usr.Email == "" {
		return usr, fmt.Errorf("could not find user with email '%s'", normalizedEmail)
	}

	return usr, nil
}

func (s *UserService) GetUserById(id uint) (m.User, error) {
	usr := m.User{
		Model: gorm.Model{ID: id},
	}

	res := s.db.Where(&usr).First(&usr)
	if res.Error != nil {
		return m.User{}, res.Error
	}

	return usr, nil
}
