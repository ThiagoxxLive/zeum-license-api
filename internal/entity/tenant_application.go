package entity

import "time"

type TenantApplication struct {
	ID            uint       `gorm:"column:id;primaryKey"`
	IDTenant      int        `gorm:"column:id_tenant"`
	IDApplication int        `gorm:"column:id_application"`
	LicenseLimit  int        `gorm:"column:license_limit"`
	LicenseUsed   int        `gorm:"column:license_used"`
	Status        bool       `gorm:"column:status"`
	Created       time.Time  `gorm:"column:created"`
	Modified      *time.Time `gorm:"column:modified"`
}

func (TenantApplication) TableName() string {
	return "tb_tenants_applications"
}
