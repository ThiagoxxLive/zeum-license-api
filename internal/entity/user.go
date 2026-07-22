package entity

import "time"

type User struct {
	ID              uint       `gorm:"column:id;primaryKey"`
	Name            string     `gorm:"column:name"`
	Email           string     `gorm:"column:email"`
	Status          bool       `gorm:"column:status"`
	LastLoginAt     *time.Time `gorm:"column:last_login_at"`
	TermsAcceptedAt *time.Time `gorm:"column:terms_accepted_at"`
}

func (User) TableName() string {
	return "tb_users"
}
