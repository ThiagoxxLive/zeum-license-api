package repository

import (
	"errors"

	"zeum-license-api/internal/entity"

	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// FindByLicenseNumber acha um tenant pelo número de licença.
// Retorna (nil, nil) quando não encontrado, análogo ao getOneOrNullResult do Doctrine.
func (r *TenantRepository) FindByLicenseNumber(licenseNumber string) (*entity.Tenant, error) {

	var tenant entity.Tenant

	result := r.db.Where("license_number = ?", licenseNumber).First(&tenant)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return &tenant, nil
}
