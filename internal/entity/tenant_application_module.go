package entity

import "time"

type TenantApplicationModule struct {
	ID                  uint       `gorm:"column:id;primaryKey"`
	IDTenantApplication int        `gorm:"column:id_tenant_application"`
	IDApplicationModule int        `gorm:"column:id_application_module"`
	Created             time.Time  `gorm:"column:created"`
	Modified            *time.Time `gorm:"column:modified"`
}

func (TenantApplicationModule) TableName() string {
	return "tb_tenants_applications_modules"
}
