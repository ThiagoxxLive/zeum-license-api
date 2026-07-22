package database

import (
	"zeum-license-api/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
}
