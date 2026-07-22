package entity

import "time"

type ApplicationModule struct {
	ID            uint       `gorm:"column:id;primaryKey"`
	IDApplication int        `gorm:"column:id_application"`
	Code          string     `gorm:"column:code"`
	Name          string     `gorm:"column:name"`
	Description   string     `gorm:"column:description"`
	Price         string     `gorm:"column:price"`
	OnDemand      bool       `gorm:"column:on_demand"`
	OnDemandPrice string     `gorm:"column:on_demand_price"`
	Created       time.Time  `gorm:"column:created"`
	Modified      *time.Time `gorm:"column:modified"`
}

func (ApplicationModule) TableName() string {
	return "tb_application_modules"
}
