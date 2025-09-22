package auth

import (
	userProviders "boilerplate_backend_api_gin/app/data/providers"
	"boilerplate_backend_api_gin/app/http/requests"
	"boilerplate_backend_api_gin/app/http/resources"
	"boilerplate_backend_api_gin/core/helpers"
	"boilerplate_backend_api_gin/core/oauth/oauth_models"
	rpProviders "boilerplate_backend_api_gin/core/roles_and_permissions/providers"
	"boilerplate_backend_api_gin/core/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(context *gin.Context) {
	var request requests.LoginRequest

	if err := validators.Validate(context, &request); err != nil {
		return
	}

	storedUser, err := userProviders.GetUserByEmail(request.Data.Attributes.Email)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"errors": []gin.H{{
			"status": "401",
			"title":  "Unauthorized",
			"detail": "The provided email or password is incorrect",
		}}})
		return
	}

	errPassword := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(request.Data.Attributes.Password))
	if errPassword != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"errors": []gin.H{{
			"status": "401",
			"title":  "Unauthorized",
			"detail": "Invalid email or password",
		}}})
		return
	}

	// TODO: Obtener clientes OAuth y permitir seleccionar uno
	clients, err := oauth_models.GetAllClients()
	if err != nil || len(clients) == 0 {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "No OAuth client available",
		}}})
		return
	}

	client := clients[0]
	token, err := oauth_models.CreateToken(storedUser.ID, client.ClientID, "")
	if err != nil {
		helpers.Logs("ERROR", "Error creating OAuth token: "+err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error generating OAuth token",
		}}})
		return
	}

	// Obtener roles del usuario
	roles, err := rpProviders.GetUserRoles(storedUser.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error getting user roles",
		}}})
		return
	}

	// Obtener permisos del usuario
	permissions, err := rpProviders.GetUserAllPermissions(storedUser.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error getting user permissions",
		}}})
		return
	}

	response := resources.NewAuthLoginResponse(storedUser, roles, permissions, token.AccessToken, token.ExpiresAt)
	context.JSON(http.StatusOK, response)
}
