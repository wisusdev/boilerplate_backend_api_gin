package auth

import (
	"net/http"
	"semita/app/data/models"
	"semita/app/data/providers"
	"semita/app/http/requests"
	"semita/app/notifications"
	"semita/config"
	"semita/core/helpers"
	"semita/core/validators"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ForgotResponseStruct struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Status  string
			Message string
		} `json:"attributes"`
	} `json:"data"`
}

func ForgotPassword(context *gin.Context) {
	var request requests.ForgotPasswordRequest

	if err := validators.Validate(context, &request); err != nil {
		return
	}

	user, err := providers.GetUserByEmail(request.Data.Attributes.Email)
	if err != nil {
		helpers.Logs("ERROR", "Error recuperando usuario: "+err.Error())
		context.JSON(http.StatusOK, gin.H{"message": "Si el email existe, se enviará un enlace de recuperación"})
		return
	}

	token := helpers.GenerateResetToken(user.Email)
	resetURL := "http://" + config.AppConfig().Url + "/auth/reset-password?token=" + token
	_ = providers.CreatePasswordReset(user.Email, token) // Guardar token en BD
	err = notifications.SendPasswordReset(user.Email, resetURL)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "No se pudo enviar el correo de recuperación",
		}}})
		return
	}

	response := ForgotResponseStruct{}
	response.Data.Type = "users"
	response.Data.Attributes.Status = "success"
	response.Data.Attributes.Message = "message.resetPasswordEmailSent"

	context.JSON(http.StatusOK, response)
}

func ResetPassword(context *gin.Context) {
	var request requests.ResetPasswordRequest

	if err := validators.Validate(context, &request); err != nil {
		return
	}
	passwordResetByToken, errorGetPasswordResetByToken := providers.GetPasswordResetByToken(request.Data.Attributes.Token)
	if errorGetPasswordResetByToken != nil {
		context.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{
			"status": "400",
			"title":  "Invalid Token",
			"detail": "Token inválido o expirado",
		}}})
		return
	}

	if time.Since(passwordResetByToken.CreatedAt) > 2*time.Hour {
		_ = providers.DeletePasswordReset(request.Data.Attributes.Token)
		context.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{
			"status": "400",
			"title":  "Token Expired",
			"detail": "Token expirado",
		}}})
		return
	}

	user, err := providers.GetUserByEmail(passwordResetByToken.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{
			"status": "400",
			"title":  "User Not Found",
			"detail": "Usuario no encontrado",
		}}})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Data.Attributes.Password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error al encriptar contraseña",
		}}})
		return
	}

	update := models.UserStruct{ID: user.ID, Password: string(hashedPassword)}
	err = providers.UpdateUserPassword(update)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "No se pudo actualizar la contraseña",
		}}})
		return
	}

	_ = providers.DeletePasswordReset(request.Data.Attributes.Token)

	var response ForgotResponseStruct = ForgotResponseStruct{}
	response.Data.Type = "users"
	response.Data.Attributes.Status = "success"
	response.Data.Attributes.Message = "message.passwordResetSuccess"

	context.JSON(http.StatusOK, response)
}
