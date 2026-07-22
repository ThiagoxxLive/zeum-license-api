package database

import (
	"time"

	"zeum-license-api/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	return db, nil
}
