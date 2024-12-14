package database

import (
	"auth/config"
	"auth/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
}

func Initialize(config config.DatabaseConfig) (Database, error) {
	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
