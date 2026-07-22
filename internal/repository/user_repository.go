package repository

import (
	"errors"

	"zeum-license-api/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByEmail acha um usuário pelo email, ignorando registros com soft-delete.
// Retorna (nil, nil) quando não encontrado, análogo ao getOneOrNullResult do Doctrine.
func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {

	var user entity.User

	result := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return &user, nil
}
