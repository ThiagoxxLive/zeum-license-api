package dto

import "time"

type ProfileApplicationResponse struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Code        string   `json:"code"`
	URL         *string  `json:"url"`
	Status      bool     `json:"status"`
	Permissions []string `json:"permissions"`
}

type ProfileTenantResponse struct {
	ID           int                          `json:"id"`
	Name         string                       `json:"name"`
	Slug         string                       `json:"slug"`
	Active       bool                         `json:"active"`
	Admin        bool                         `json:"admin"`
	Applications []ProfileApplicationResponse `json:"applications"`
}

type ProfileResponse struct {
	ID              int                     `json:"id"`
	Name            string                  `json:"name"`
	Email           string                  `json:"email"`
	Active          bool                    `json:"active"`
	LastLoginAt     *time.Time              `json:"last_login_at"`
	TermsAcceptedAt *time.Time              `json:"terms_accepted_at"`
	Tenants         []ProfileTenantResponse `json:"tenants"`
}
