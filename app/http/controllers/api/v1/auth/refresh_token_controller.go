package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"semita/core/oauth/oauth_models"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
		ClientID     string `json:"client_id" binding:"required"`
		ClientSecret string `json:"client_secret" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	// Validar credenciales del cliente
	_, err := oauth_models.ValidateClientCredentials(request.ClientID, request.ClientSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales de cliente inválidas"})
		return
	}

	// Renovar token
	token, err := oauth_models.RefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token inválido"})
		return
	}

	// Devolver el nuevo token
	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		RefreshToken: token.RefreshToken,
		Scope:        token.Scopes,
	})
}
