// Package middleware contém os middlewares HTTP do serviço, incluindo a
// autenticação por X-API-Key equivalente a App\Security\ApiKeyAuthenticator
// do zeum-admin-api.
package middleware

import (
	"crypto/subtle"
	"log"
	"net/http"

	"zeum-license-api/internal/config"
	"zeum-license-api/internal/cryptox"
	"zeum-license-api/internal/repository"

	"github.com/gin-gonic/gin"
)

const apiKeyHeaderName = "X-API-Key"

const applicationContextKey = "application"

// APIKeyAuth exige exclusivamente o header X-API-Key, sem o fallback para
// Bearer JWT do Keycloak feito por Auth(). Use em rotas que só devem ser
// chamadas por integrações server-to-server autenticadas via API Key.
func APIKeyAuth(cfg *config.Config, applicationRepository *repository.ApplicationRepository) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		apiKeyAuth(ctx, cfg, applicationRepository)
	}
}

// apiKeyAuth valida o header X-API-Key contra as Applications ativas,
// equivalente a App\Security\ApiKeyAuthenticator::authenticate().
func apiKeyAuth(ctx *gin.Context, cfg *config.Config, applicationRepository *repository.ApplicationRepository) {

	apiKey := ctx.GetHeader(apiKeyHeaderName)

	if apiKey == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": true, "message": "API Key não informada."})
		return
	}

	rawKey, err := cryptox.ParseKey(cfg.DefuseSecret)

	if err != nil {
		log.Printf("DEFUSE_SECRET inválido: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Ocorreu um erro interno. Tente novamente mais tarde."})
		return
	}

	applications, err := applicationRepository.FindAllWithAPIKey()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}

	for _, application := range applications {

		if application.APIKey == nil {
			continue
		}

		decrypted, err := cryptox.Decrypt(*application.APIKey, rawKey)

		if err != nil {
			continue
		}

		if subtle.ConstantTimeCompare([]byte(decrypted), []byte(apiKey)) == 1 {
			ctx.Set(applicationContextKey, application)
			ctx.Next()
			return
		}
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": true, "message": "API Key inválida."})
}
