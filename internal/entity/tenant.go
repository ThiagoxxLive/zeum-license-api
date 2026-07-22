package entity

import "time"

type Tenant struct {
	ID               uint       `gorm:"column:id;primaryKey"`
	IDCustomer       int        `gorm:"column:id_customer"`
	Name             string     `gorm:"column:name"`
	Slug             string     `gorm:"column:slug"`
	LicenseNumber    string     `gorm:"column:license_number"`
	LicenseExpiresAt *time.Time `gorm:"column:license_expires_at"`
	Status           bool       `gorm:"column:status"`
	Created          time.Time  `gorm:"column:created"`
	Modified         *time.Time `gorm:"column:modified"`
}

func (Tenant) TableName() string {
	return "tb_tenants"
}
