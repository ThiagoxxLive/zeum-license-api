package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"zeum-license-api/internal/config"
	"zeum-license-api/internal/controller"
	"zeum-license-api/internal/database"
	"zeum-license-api/internal/keycloak"
	"zeum-license-api/internal/repository"
	"zeum-license-api/internal/router"
	"zeum-license-api/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewKeycloakPublicKeyProvider(cfg *config.Config) *keycloak.PublicKeyProvider {
	return keycloak.NewPublicKeyProvider(cfg.KeycloakBaseURL, cfg.KeycloakRealm)
}

func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config, engine *gin.Engine) *http.Server {

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("erro ao iniciar o servidor: %v", err)
				}
			}()

			log.Printf("servidor ouvindo em %s", server.Addr)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}

func main() {

	fx.New(
		fx.Provide(
			config.Load,
			database.NewConnection,
			repository.NewApplicationRepository,
			repository.NewTenantRepository,
			repository.NewTenantApplicationRepository,
			repository.NewTenantApplicationModuleRepository,
			repository.NewUserRepository,
			repository.NewTenantUserRepository,
			repository.NewUserRoleRepository,
			NewKeycloakPublicKeyProvider,
			service.NewLicenseService,
			service.NewProfileService,
			controller.NewLicenseController,
			controller.NewProfileController,
			controller.NewHealthController,
			router.NewRouter,
			NewHTTPServer,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
