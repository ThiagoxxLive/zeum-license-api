package router

import (
	"time"

	"zeum-license-api/internal/config"
	"zeum-license-api/internal/controller"
	"zeum-license-api/internal/keycloak"
	"zeum-license-api/internal/middleware"
	"zeum-license-api/internal/ratelimit"
	"zeum-license-api/internal/repository"
	"zeum-license-api/internal/service"

	"github.com/gin-gonic/gin"
)

// Mesmos limites do RateLimitSubscriber do zeum-admin-api (api_global/api_strict).
const (
	globalRateLimit = 120
	strictRateLimit = 20
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

	globalLimiter := ratelimit.NewSlidingWindow(globalRateLimit, time.Minute)
	strictLimiter := ratelimit.NewSlidingWindow(strictRateLimit, time.Minute)

	v1 := engine.Group("/api/v1")
	v1.Use(middleware.Auth(cfg, applicationRepository, keyProvider))
	v1.POST("/license", middleware.RateLimit(globalLimiter), licenseController.FindAction)
	v1.GET(
		"/profile/:email",
		middleware.RateLimit(strictLimiter),
		middleware.RequireOwnerOrAdmin(profileService),
		profileController.FindAction,
	)

	return engine
}
