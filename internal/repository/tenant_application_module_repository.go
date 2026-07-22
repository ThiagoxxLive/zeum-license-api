package repository

import (
	"time"

	"gorm.io/gorm"
)

type TenantApplicationModuleForLicense struct {
	ID                  int
	IDTenantApplication int
	IDApplicationModule int
	Code                string
	Name                string
	Description         string
	Price               string
	OnDemand            bool
	OnDemandPrice       string
	Created             time.Time
	Modified            *time.Time
}

type TenantApplicationModuleRepository struct {
	db *gorm.DB
}

func NewTenantApplicationModuleRepository(db *gorm.DB) *TenantApplicationModuleRepository {
	return &TenantApplicationModuleRepository{db: db}
}

// FindByTenantApplicationID acha os módulos vinculados a uma aplicação de tenant, para uso na licença.
func (r *TenantApplicationModuleRepository) FindByTenantApplicationID(tenantApplicationID int) ([]TenantApplicationModuleForLicense, error) {

	var results []TenantApplicationModuleForLicense

	err := r.db.Table("tb_tenants_applications_modules AS tam").
		Select(`
			tam.id AS id,
			tam.id_tenant_application AS id_tenant_application,
			tam.id_application_module AS id_application_module,
			am.code AS code,
			am.name AS name,
			am.description AS description,
			am.price AS price,
			am.on_demand AS on_demand,
			am.on_demand_price AS on_demand_price,
			tam.created AS created,
			tam.modified AS modified
		`).
		Joins("INNER JOIN tb_application_modules AS am ON tam.id_application_module = am.id").
		Where("tam.id_tenant_application = ?", tenantApplicationID).
		Order("tam.id ASC").
		Scan(&results).Error

	return results, err
}
