package auth

import (
	userModels "boilerplate_backend_api_gin/app/data/models"
	userProviders "boilerplate_backend_api_gin/app/data/providers"
	"boilerplate_backend_api_gin/app/http/requests"
	"boilerplate_backend_api_gin/core/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterResponse struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Status  string
			Message string
		} `json:"attributes"`
	} `json:"data"`
}

func Register(context *gin.Context) {
	var request requests.RegisterRequest

	if err := validators.Validate(context, &request); err != nil {
		return
	}

	existingUser, _ := userProviders.GetUserByEmail(request.Data.Attributes.Email)
	if existingUser.ID != "" {
		context.JSON(http.StatusConflict, gin.H{"errors": []gin.H{{
			"status": "409",
			"title":  "Conflict",
			"detail": "Email already registered",
		}}})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Data.Attributes.Password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error encrypting password",
		}}})
		return
	}

	userToStore := userModels.UserStruct{
		Username:  request.Data.Attributes.Username,
		FirstName: request.Data.Attributes.Firstname,
		LastName:  request.Data.Attributes.Lastname,
		Email:     request.Data.Attributes.Email,
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

	response := RegisterResponse{}
	response.Data.Type = "users"
	response.Data.Attributes.Status = "success"
	response.Data.Attributes.Message = "User registered successfully"

	context.JSON(http.StatusCreated, response)
}
