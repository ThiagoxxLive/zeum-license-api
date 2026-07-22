package middleware

import (
	"errors"
	"net/http"
	"strings"

	"zeum-license-api/internal/keycloak"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const authorizationHeaderName = "Authorization"

const userIDClaimContextKey = "keycloak_user_id"

const userEmailClaimContextKey = "keycloak_user_email"

// jwtFailureJSON reproduz o formato de Lexik\...\JWTAuthenticationFailureResponse:
// {"code":401,"message":"..."} com o header WWW-Authenticate: Bearer.
func jwtFailureJSON(ctx *gin.Context, message string) {
	ctx.Header("WWW-Authenticate", "Bearer")
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": message})
}

// keycloakJWTAuth valida um Bearer token emitido pelo Keycloak, equivalente ao
// KeycloakJWSProvider + lexik JWTAuthenticator do zeum-admin-api.
func keycloakJWTAuth(ctx *gin.Context, keyProvider *keycloak.PublicKeyProvider) {

	tokenString := extractBearerToken(ctx)

	token, err := parseAndVerify(tokenString, keyProvider, false)

	if err != nil {
		// mesmo retry do PHP: a chave em cache pode ter rotacionado no Keycloak.
		token, err = parseAndVerify(tokenString, keyProvider, true)
	}

	if err != nil {

		if errors.Is(err, jwt.ErrTokenExpired) {
			jwtFailureJSON(ctx, "Expired JWT Token")
			return
		}

		jwtFailureJSON(ctx, "Invalid JWT Token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		jwtFailureJSON(ctx, "Invalid JWT Token")
		return
	}

	sub, ok := claims["sub"]

	if !ok {
		jwtFailureJSON(ctx, `Unable to find key "sub" in the token payload.`)
		return
	}

	ctx.Set(userIDClaimContextKey, sub)

	if email, ok := claims["email"].(string); ok {
		ctx.Set(userEmailClaimContextKey, email)
	}

	ctx.Next()
}

func extractBearerToken(ctx *gin.Context) string {

	header := ctx.GetHeader(authorizationHeaderName)

	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}

	return strings.TrimPrefix(header, "Bearer ")
}

func parseAndVerify(tokenString string, keyProvider *keycloak.PublicKeyProvider, forceRefresh bool) (*jwt.Token, error) {

	publicKey, err := keyProvider.PublicKey(forceRefresh)

	if err != nil {
		return nil, err
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return publicKey, nil

	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
}
