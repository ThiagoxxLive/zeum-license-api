package middleware

import (
	"net/http"

	"zeum-license-api/internal/service"

	"github.com/gin-gonic/gin"
)

// RequireOwnerOrAdmin reproduz o OwnerOrAdminListener do zeum-admin-api: requisições
// autenticadas por X-API-Key sempre passam; requisições autenticadas por Bearer JWT só
// passam se o email do token bater com o :email da rota, ou se o usuário do token for
// admin em algum tenant.
func RequireOwnerOrAdmin(profileService *service.ProfileService) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		if _, ok := ctx.Get(applicationContextKey); ok {
			ctx.Next()
			return
		}

		tokenEmail, ok := ctx.Get(userEmailClaimContextKey)

		if !ok {
			denyOwnerOrAdmin(ctx)
			return
		}

		email, ok := tokenEmail.(string)

		if !ok || email == "" {
			denyOwnerOrAdmin(ctx)
			return
		}

		if email == ctx.Param("email") {
			ctx.Next()
			return
		}

		profile, err := profileService.FindByEmail(email)

		if err != nil {
			denyOwnerOrAdmin(ctx)
			return
		}

		for _, tenant := range profile.Tenants {

			if tenant.Admin {
				ctx.Next()
				return
			}
		}

		denyOwnerOrAdmin(ctx)
	}
}

func denyOwnerOrAdmin(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": true, "message": "Acesso negado."})
}
