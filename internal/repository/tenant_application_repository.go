package repository

import "gorm.io/gorm"

type TenantApplicationForLicense struct {
	TenantApplicationID int
	ID                   int
	Name                 string
	Description          *string
	Code                 string
	URL                  *string
	Status               bool
}

type TenantApplicationForTenants struct {
	IDTenant    int
	ID          int
	Name        string
	Description *string
	Code        string
	URL         *string
	Status      bool
}

type TenantApplicationRepository struct {
	db *gorm.DB
}

func NewTenantApplicationRepository(db *gorm.DB) *TenantApplicationRepository {
	return &TenantApplicationRepository{db: db}
}

// FindByTenantIDForLicense acha as aplicações vinculadas a um tenant, para uso na licença.
func (r *TenantApplicationRepository) FindByTenantIDForLicense(tenantID uint) ([]TenantApplicationForLicense, error) {

	var results []TenantApplicationForLicense

	err := r.db.Table("tb_tenants_applications AS ta").
		Select(`
			ta.id AS tenant_application_id,
			a.id AS id,
			a.name AS name,
			a.description AS description,
			a.slug AS code,
			a.url AS url,
			ta.status AS status
		`).
		Joins("INNER JOIN tb_applications AS a ON ta.id_application = a.id").
		Where("ta.id_tenant = ?", tenantID).
		Scan(&results).Error

	return results, err
}

// FindByTenantIDs acha as aplicações vinculadas a uma lista de tenants, para uso no profile.
func (r *TenantApplicationRepository) FindByTenantIDs(tenantIDs []int) ([]TenantApplicationForTenants, error) {

	var results []TenantApplicationForTenants

	err := r.db.Table("tb_tenants_applications AS ta").
		Select(`
			ta.id_tenant AS id_tenant,
			a.id AS id,
			a.name AS name,
			a.description AS description,
			a.slug AS code,
			a.url AS url,
			ta.status AS status
		`).
		Joins("INNER JOIN tb_applications AS a ON ta.id_application = a.id").
		Where("ta.id_tenant IN ?", tenantIDs).
		Scan(&results).Error

	return results, err
}
