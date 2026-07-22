package repository

import "gorm.io/gorm"

type TenantUserTenant struct {
	ID     int
	Name   string
	Slug   string
	Active bool
	Admin  bool
}

type TenantUserRepository struct {
	db *gorm.DB
}

func NewTenantUserRepository(db *gorm.DB) *TenantUserRepository {
	return &TenantUserRepository{db: db}
}

// FindTenantsByUserID acha os tenants aos quais um usuário está vinculado.
func (r *TenantUserRepository) FindTenantsByUserID(userID uint) ([]TenantUserTenant, error) {

	var results []TenantUserTenant

	err := r.db.Table("tb_tenants_users AS tu").
		Select(`
			t.id AS id,
			t.name AS name,
			t.slug AS slug,
			t.status AS active,
			tu.admin AS admin
		`).
		Joins("INNER JOIN tb_tenants AS t ON tu.id_tenant = t.id").
		Where("tu.id_user = ? AND tu.deleted_at IS NULL", userID).
		Scan(&results).Error

	return results, err
}
