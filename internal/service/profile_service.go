package service

import (
	"zeum-license-api/internal/apperror"
	"zeum-license-api/internal/dto"
	"zeum-license-api/internal/repository"
)

type ProfileService struct {
	userRepository              *repository.UserRepository
	tenantUserRepository        *repository.TenantUserRepository
	tenantApplicationRepository *repository.TenantApplicationRepository
	userRoleRepository          *repository.UserRoleRepository
}

func NewProfileService(
	userRepository *repository.UserRepository,
	tenantUserRepository *repository.TenantUserRepository,
	tenantApplicationRepository *repository.TenantApplicationRepository,
	userRoleRepository *repository.UserRoleRepository,
) *ProfileService {

	return &ProfileService{
		userRepository:              userRepository,
		tenantUserRepository:        tenantUserRepository,
		tenantApplicationRepository: tenantApplicationRepository,
		userRoleRepository:          userRoleRepository,
	}
}

func (s *ProfileService) FindByEmail(email string) (*dto.ProfileResponse, error) {

	user, err := s.userRepository.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.ErrProfileNotFound
	}

	response := &dto.ProfileResponse{
		ID:              int(user.ID),
		Name:            user.Name,
		Email:           user.Email,
		Active:          user.Status,
		LastLoginAt:     user.LastLoginAt,
		TermsAcceptedAt: user.TermsAcceptedAt,
		Tenants:         []dto.ProfileTenantResponse{},
	}

	tenants, err := s.tenantUserRepository.FindTenantsByUserID(user.ID)

	if err != nil {
		return nil, err
	}

	if len(tenants) == 0 {
		return response, nil
	}

	tenantIDs := make([]int, 0, len(tenants))

	for _, tenant := range tenants {
		tenantIDs = append(tenantIDs, tenant.ID)
	}

	applications, err := s.tenantApplicationRepository.FindByTenantIDs(tenantIDs)

	if err != nil {
		return nil, err
	}

	permissionsByApplication, err := s.userRoleRepository.FindPermissionCodesByUserID(user.ID)

	if err != nil {
		return nil, err
	}

	applicationsByTenant := make(map[int][]dto.ProfileApplicationResponse)

	for _, application := range applications {

		applicationsByTenant[application.IDTenant] = append(applicationsByTenant[application.IDTenant], dto.ProfileApplicationResponse{
			ID:          application.ID,
			Name:        application.Name,
			Description: application.Description,
			Code:        application.Code,
			URL:         application.URL,
			Status:      application.Status,
			Permissions: permissionsByApplication[application.ID],
		})
	}

	for _, tenant := range tenants {

		response.Tenants = append(response.Tenants, dto.ProfileTenantResponse{
			ID:           tenant.ID,
			Name:         tenant.Name,
			Slug:         tenant.Slug,
			Active:       tenant.Active,
			Admin:        tenant.Admin,
			Applications: applicationsByTenant[tenant.ID],
		})
	}

	return response, nil
}
