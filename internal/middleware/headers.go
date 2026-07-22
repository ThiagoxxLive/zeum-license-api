package middleware

import (
	"runtime"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders replica os headers de resposta expostos pelo zeum-admin-api
// (CORS liberado + hardening básico), para manter paridade entre as APIs.
func SecurityHeaders() gin.HandlerFunc {
	poweredBy := "Go/" + runtime.Version()[2:]

	return func(ctx *gin.Context) {
		ctx.Header("X-Powered-By", poweredBy)
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "X-Cache")
		ctx.Header("X-Robots-Tag", "noindex")
		ctx.Header("X-Frame-Options", "SAMEORIGIN")
		ctx.Header("X-Content-Type-Options", "nosniff")

		ctx.Next()
	}
}
