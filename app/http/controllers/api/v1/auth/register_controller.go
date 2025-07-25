package auth

import (
	"net/http"
	"semita/app/data/models"
	"semita/app/data/structs"
	"semita/app/http/requests"
	"semita/app/http/resources"
	"semita/core/helpers"
	"semita/core/oauth/oauth_models"

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

	existingUser, _ := models.GetUserByEmail(req.Email)
	if existingUser.ID > 0 {
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

	userToStore := structs.StoreUserStruct{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
	}

	errorStore := models.StoreUser(userToStore)
	if errorStore != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error saving user to the database",
		}}})
		return
	}

	storedUser, err := models.GetUserByEmail(req.Email)
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
	token, err := oauth_models.CreateToken(int64(storedUser.ID), client.ID, "")
	if err != nil {
		helpers.Logs("ERROR", "Error generating OAuth token: "+err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "User created but error generating OAuth token",
		}}})
		return
	}

	resource := resources.NewAuthResource(uint(storedUser.ID), storedUser.FirstName+" "+storedUser.LastName, storedUser.Email, token.AccessToken)
	context.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"type": "users",
			"id":   resource.ID,
			"attributes": gin.H{
				"name":  resource.Name,
				"email": resource.Email,
			},
			"meta": gin.H{
				"token":         resource.Token,
				"refresh_token": token.RefreshToken,
				"expires_in":    86400,
				"scope":         token.Scopes,
			},
		},
	})
}
