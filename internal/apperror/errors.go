package apperror

import "errors"

var (
	ErrLicenseNotFound = errors.New("Licença não encontrada.")
	ErrLicenseExpired  = errors.New("Licença expirada.")
	ErrProfileNotFound = errors.New("Perfil do usuário não encontrado.")
)
