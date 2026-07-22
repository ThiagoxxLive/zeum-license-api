package middleware

import (
	"net/http"

	"zeum-license-api/internal/config"
	"zeum-license-api/internal/keycloak"
	"zeum-license-api/internal/repository"

	"github.com/gin-gonic/gin"
)

// Auth reproduz o firewall "api" (pattern ^/v1) do zeum-admin-api: aceita tanto
// um Bearer JWT emitido pelo Keycloak quanto um header X-API-Key. Quando nenhum
// dos dois é enviado, a requisição nunca chega a "autenticar" em nenhum dos dois
// authenticators do Symfony e cai na regra de access_control (IS_AUTHENTICATED_FULLY),
// que resulta em 403 "Access Denied. The user is not appropriately authenticated."
func Auth(
	cfg *config.Config,
	applicationRepository *repository.ApplicationRepository,
	keyProvider *keycloak.PublicKeyProvider,
) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		if extractBearerToken(ctx) != "" {
			keycloakJWTAuth(ctx, keyProvider)
			return
		}

		if ctx.Request.Header.Values(apiKeyHeaderName) != nil {
			apiKeyAuth(ctx, cfg, applicationRepository)
			return
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "Access Denied. The user is not appropriately authenticated.",
		})
	}
}
