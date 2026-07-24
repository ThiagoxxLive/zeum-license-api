package service

import (
	"time"
	"zeum-license-api/internal/apperror"
	"zeum-license-api/internal/dto"
	"zeum-license-api/internal/repository"
)

type LicenseService struct {
	tenantRepository                  *repository.TenantRepository
	tenantApplicationRepository       *repository.TenantApplicationRepository
	tenantApplicationModuleRepository *repository.TenantApplicationModuleRepository
}

func NewLicenseService(
	tenantRepository *repository.TenantRepository,
	tenantApplicationRepository *repository.TenantApplicationRepository,
	tenantApplicationModuleRepository *repository.TenantApplicationModuleRepository,
) *LicenseService {

	return &LicenseService{
		tenantRepository:                  tenantRepository,
		tenantApplicationRepository:       tenantApplicationRepository,
		tenantApplicationModuleRepository: tenantApplicationModuleRepository,
	}
}

func (s *LicenseService) FindByLicenseNumber(licenseNumber string) (*dto.LicenseResponse, error) {

	tenant, err := s.tenantRepository.FindByLicenseNumber(licenseNumber)

	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, apperror.ErrLicenseNotFound
	}

	if tenant.LicenseExpiresAt != nil && tenant.LicenseExpiresAt.Before(time.Now()) {
		return nil, apperror.ErrLicenseExpired
	}

	tenantApplications, err := s.tenantApplicationRepository.FindByTenantIDForLicense(tenant.ID)

	if err != nil {
		return nil, err
	}

	applications := make([]dto.ApplicationResponse, 0, len(tenantApplications))

	for _, tenantApplication := range tenantApplications {

		modules, err := s.tenantApplicationModuleRepository.FindByTenantApplicationID(tenantApplication.TenantApplicationID)

		if err != nil {
			return nil, err
		}

		moduleResponses := make([]dto.ModuleResponse, 0, len(modules))

		for _, module := range modules {
			moduleResponses = append(moduleResponses, dto.ModuleResponse{
				ID:                  module.ID,
				IDTenantApplication: module.IDTenantApplication,
				IDApplicationModule: module.IDApplicationModule,
				Code:                module.Code,
				Name:                module.Name,
				Description:         module.Description,
				Price:               module.Price,
				OnDemand:            module.OnDemand,
				OnDemandPrice:       module.OnDemandPrice,
				Created:             module.Created,
				Modified:            module.Modified,
			})
		}

		applications = append(applications, dto.ApplicationResponse{
			ID:          tenantApplication.ID,
			Name:        tenantApplication.Name,
			Description: tenantApplication.Description,
			Code:        tenantApplication.Code,
			URL:         tenantApplication.URL,
			Status:      tenantApplication.Status,
			Modules:     moduleResponses,
		})
	}

	return &dto.LicenseResponse{
		Tenant: dto.TenantResponse{
			ID:               int(tenant.ID),
			Name:             tenant.Name,
			Slug:             tenant.Slug,
			LicenseNumber:    tenant.LicenseNumber,
			LicenseExpiresAt: tenant.LicenseExpiresAt,
			Status:           tenant.Status,
		},
		Applications: applications,
	}, nil
}

func (s *LicenseService) FindLicenseNumberByTenantID(tenantID uint) (*dto.LicenseNumberResponse, error) {

	tenant, err := s.tenantRepository.FindByID(tenantID)

	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, apperror.ErrTenantNotFound
	}

	return &dto.LicenseNumberResponse{
		LicenseNumber: tenant.LicenseNumber,
	}, nil
}
