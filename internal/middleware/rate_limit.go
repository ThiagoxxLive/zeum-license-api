package middleware

import (
	"net/http"
	"strconv"

	"zeum-license-api/internal/ratelimit"

	"github.com/gin-gonic/gin"
)

// RateLimit reproduz o RateLimitSubscriber do zeum-admin-api: janela deslizante
// por IP, respondendo 429 com Retry-After quando o limite estoura.
func RateLimit(limiter *ratelimit.SlidingWindow) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		allowed, retryAfter := limiter.Allow(ctx.ClientIP())

		if !allowed {

			seconds := int(retryAfter.Seconds())

			if seconds < 1 {
				seconds = 1
			}

			ctx.Header("Retry-After", strconv.Itoa(seconds))
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Muitas requisições. Tente novamente em instantes."})
			return
		}

		ctx.Next()
	}
}
