package dto

import "time"

type FindLicenseRequest struct {
	LicenseNumber string `json:"license_number" binding:"required,max=100"`
}

type TenantResponse struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	LicenseNumber    string     `json:"license_number"`
	LicenseExpiresAt *time.Time `json:"license_expires_at"`
	Status           bool       `json:"status"`
}

type ModuleResponse struct {
	ID                  int        `json:"id"`
	IDTenantApplication int        `json:"id_tenant_application"`
	IDApplicationModule int        `json:"id_application_module"`
	Code                string     `json:"code"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Price               string     `json:"price"`
	OnDemand            bool       `json:"on_demand"`
	OnDemandPrice       string     `json:"on_demand_price"`
	Created             time.Time  `json:"created"`
	Modified            *time.Time `json:"modified"`
}

type ApplicationResponse struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Code        string           `json:"code"`
	URL         *string          `json:"url"`
	Status      bool             `json:"status"`
	Modules     []ModuleResponse `json:"modules"`
}

type LicenseResponse struct {
	Tenant       TenantResponse         `json:"tenant"`
	Applications []ApplicationResponse `json:"applications"`
}

type LicenseNumberResponse struct {
	LicenseNumber string `json:"license_number"`
}
