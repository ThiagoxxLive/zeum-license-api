package entity

type Application struct {
	ID          uint    `gorm:"column:id;primaryKey"`
	Name        string  `gorm:"column:name"`
	Slug        string  `gorm:"column:slug"`
	Description *string `gorm:"column:description"`
	URL         *string `gorm:"column:url"`
	Status      bool    `gorm:"column:status"`
	Internal    bool    `gorm:"column:internal"`
	HasModules  bool    `gorm:"column:has_modules"`
	APIKey      *string `gorm:"column:api_key"`
}

func (Application) TableName() string {
	return "tb_applications"
}
