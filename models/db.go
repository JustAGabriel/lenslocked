package models

import (
	"fmt"

	"github.com/pressly/goose"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func (pc PostgresConfig) String() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", pc.Host, pc.User, pc.Password, pc.Database, pc.Port, pc.SSLMode)
}

func GetDefaultDBConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
}

func Migrate(db *gorm.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migration, set dialect: %w", err)
	}

	standardDB, err2 := db.DB()
	if err2 != nil {
		return fmt.Errorf("migration, db conversion: %w", err2)
	}

	err = goose.Up(standardDB, dir)
	if err != nil {
		return fmt.Errorf("migration, up: %w", err)
	}

	return nil
}
