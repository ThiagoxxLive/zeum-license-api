package router

import (
	"zeum-license-api/internal/config"
	"zeum-license-api/internal/controller"
	"zeum-license-api/internal/keycloak"
	"zeum-license-api/internal/middleware"
	"zeum-license-api/internal/repository"
	"zeum-license-api/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	applicationRepository *repository.ApplicationRepository,
	keyProvider *keycloak.PublicKeyProvider,
	licenseController *controller.LicenseController,
	profileController *controller.ProfileController,
	profileService *service.ProfileService,
) *gin.Engine {

	engine := gin.Default()
	engine.Use(middleware.SecurityHeaders())

	v1 := engine.Group("/v1")
	v1.Use(middleware.Auth(cfg, applicationRepository, keyProvider))
	v1.POST("/license", licenseController.FindAction)
	v1.GET("/profile/:email", middleware.RequireOwnerOrAdmin(profileService), profileController.FindAction)

	return engine
}
