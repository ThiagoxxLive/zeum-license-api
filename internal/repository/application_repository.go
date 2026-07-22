package repository

import (
	"zeum-license-api/internal/entity"

	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

// FindAllWithAPIKey acha as aplicações ativas que possuem API Key configurada,
// para uso na autenticação por X-API-Key.
func (r *ApplicationRepository) FindAllWithAPIKey() ([]entity.Application, error) {

	var applications []entity.Application

	err := r.db.
		Where("api_key IS NOT NULL").
		Where("status = ?", true).
		Find(&applications).Error

	return applications, err
}
