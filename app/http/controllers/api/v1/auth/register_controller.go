package auth

import (
	"net/http"
	userModels "semita/app/data/models"
	userProviders "semita/app/data/providers"
	"semita/app/http/requests"
	"semita/app/http/resources"
	"semita/core/helpers"
	"semita/core/oauth/oauth_models"
	rpProviders "semita/core/roles_and_permissions/providers"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(context *gin.Context) {
	var req requests.RegisterRequest
	if err := req.Validate(context); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{
			"status": "400",
			"title":  "Validation Error",
			"detail": err.Error(),
		}}})
		return
	}

	existingUser, _ := userProviders.GetUserByEmail(req.Email)
	if existingUser.ID != "" {
		context.JSON(http.StatusConflict, gin.H{"errors": []gin.H{{
			"status": "409",
			"title":  "Conflict",
			"detail": "Email already registered",
		}}})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error encrypting password",
		}}})
		return
	}

	userToStore := userModels.UserStruct{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
	}

	_, errorStore := userProviders.StoreUser(userToStore)
	if errorStore != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error saving user to the database",
		}}})
		return
	}

	storedUser, err := userProviders.GetUserByEmail(req.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "User created but error retrieving user data",
		}}})
		return
	}

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
		helpers.Logs("ERROR", "Error generating OAuth token: "+err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "User created but error generating OAuth token",
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
	context.JSON(http.StatusCreated, response)
}
