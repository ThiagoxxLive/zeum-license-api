package repository

import "gorm.io/gorm"

type applicationPermission struct {
	IDApplication int
	Code          string
}

type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

// FindPermissionCodesByUserID retorna os códigos de permissão do usuário,
// agrupados por aplicação (id_application), com base nos roles que ele possui.
func (r *UserRoleRepository) FindPermissionCodesByUserID(userID uint) (map[int][]string, error) {

	var rows []applicationPermission

	err := r.db.Table("tb_user_roles AS ur").
		Select("DISTINCT ap.id_application AS id_application, ap.code AS code").
		Joins("INNER JOIN tb_roles AS r ON r.id = ur.id_role").
		Joins("INNER JOIN tb_roles_permissions AS rp ON rp.id_role = r.id").
		Joins("INNER JOIN tb_applications_permissions AS ap ON ap.id = rp.id_application_permission").
		Where("ur.id_user = ?", userID).
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	grouped := make(map[int][]string)

	for _, row := range rows {
		grouped[row.IDApplication] = append(grouped[row.IDApplication], row.Code)
	}

	return grouped, nil
}
