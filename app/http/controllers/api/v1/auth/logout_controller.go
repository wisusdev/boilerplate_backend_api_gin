package auth

import (
	"boilerplate_backend_api_gin/core/oauth/oauth_models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(context *gin.Context) {
	var tokenObj, exists = context.Get("token")

	if !exists {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Token no encontrado en el contexto",
		})
		return
	}

	// Revoke the token in the database
	var token, ok = tokenObj.(*oauth_models.OAuthToken)
	if !ok {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Token no es una cadena válida",
		})
		return
	}

	tokenString := token.AccessToken

	err := oauth_models.RevokeToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al revocar el token: " + err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Sesión cerrada correctamente",
	})
}
