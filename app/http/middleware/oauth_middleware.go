package middleware

import (
	"fmt"
	"net/http"
	"semita/core/helpers"
	"semita/core/oauth/oauth_models"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware es el middleware de autenticación OAuth2
func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")

		if authHeader == "" {
			helpers.Logs("ERROR", "Authorization header is required")
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token no proporcionado",
			})
			return
		}

		acceptTypeHeader := "application/vnd.api+json"

		acceptHeader := context.GetHeader("Accept")
		contentTypeHeader := context.GetHeader("Content-Type")

		if !strings.Contains(acceptHeader, acceptTypeHeader) && !strings.Contains(contentTypeHeader, acceptTypeHeader) {
			helpers.Logs("ERROR", fmt.Sprintf("Accept or Content-Type header must be %s", acceptTypeHeader))
			context.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"error": fmt.Sprintf("El encabezado Accept o Content-Type debe ser %s", acceptTypeHeader),
			})
			return
		}

		// El token debe tener el formato "Bearer {token}"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.Logs("ERROR", "Invalid token format")
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Formato de token inválido",
			})
			return
		}

		tokenString := parts[1]

		// Debug: Log token info
		tokenStart := tokenString
		if len(tokenString) > 50 {
			tokenStart = tokenString[:50]
		}
		helpers.Logs("DEBUG", fmt.Sprintf("Received token length: %d, starts with: %s", len(tokenString), tokenStart))

		// Validar el token JWT
		claims, err := helpers.ValidateJWTToken(tokenString)
		if err != nil {
			helpers.Logs("ERROR", "Invalid token: "+err.Error())
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido",
			})
			return
		}

		// Verificar si el token existe en la base de datos y no está revocado
		// Usamos GetTokenByJTI para evitar doble validación del JWT
		token, err := oauth_models.GetTokenByJTI(claims.JTI, tokenString)

		if err != nil || token.Revoked {
			helpers.Logs("ERROR", "Token revoked or not found in DB: "+err.Error())
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token revocado o inválido",
			})
			return
		}

		// Almacenar información del token para uso posterior en controladores
		context.Set("user_id", claims.Subject)
		context.Set("client_id", claims.Audience)
		context.Set("token_id", claims.JTI)
		context.Set("token_scopes", claims.Scopes)
		context.Set("token", token)

		context.Next()
	}
}

// ScopeMiddleware es el middleware para verificar los scopes requeridos
func ScopeMiddleware(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Este middleware debe usarse después de AuthMiddleware
		scopes, exists := c.Get("token_scopes")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No se ha autenticado correctamente",
			})
			return
		}

		tokenScopes, ok := scopes.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error al procesar los scopes",
			})
			return
		}

		// Verificar si el token tiene al menos uno de los scopes requeridos
		for _, requiredScope := range requiredScopes {
			if helpers.HasScope(tokenScopes, requiredScope) {
				// Sí tiene al menos uno de los scopes requeridos, permitir el acceso
				c.Next()
				return
			}
		}

		// Si no tiene ninguno de los scopes requeridos, denegar el acceso
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":           "Acceso denegado",
			"required_scopes": requiredScopes,
		})
	}
}
