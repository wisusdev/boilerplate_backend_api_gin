package auth

import (
	"net/http"
	"semita/app/data/models"
	"semita/app/http/requests"
	"semita/app/http/resources"
	"semita/core/oauth/oauth_models"
	"semita/core/validators"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(context *gin.Context) {
	var request requests.LoginRequest

	if err := validators.Validate(context, &request); err != nil {
		return
	}

	storedUser, err := models.GetUserByEmail(request.Data.Attributes.Email)
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
	token, err := oauth_models.CreateToken(int64(storedUser.ID), client.ID, "")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error generating OAuth token",
		}}})
		return
	}

	resource := resources.NewAuthResource(uint(storedUser.ID), storedUser.FirstName+" "+storedUser.LastName, storedUser.Email, token.AccessToken)
	response := resources.NewAuthLoginResponse(resource, token.RefreshToken, 86400, token.Scopes)
	context.JSON(http.StatusOK, response)
}
