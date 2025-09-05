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

func ForgotPassword(context *gin.Context) {
	var req requests.ForgotPasswordRequest

	if err := validators.Validate(context, &req); err != nil {
		return
	}

	user, err := providers.GetUserByEmail(req.Email)
	if err != nil {
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
	context.JSON(http.StatusOK, gin.H{"message": "Si el email existe, se enviará un enlace de recuperación"})
}

func ResetPassword(context *gin.Context) {
	var req requests.ResetPasswordRequest
	if err := validators.Validate(context, &req); err != nil {
		return
	}

	pr, err := providers.GetPasswordResetByToken(req.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{
			"status": "400",
			"title":  "Invalid Token",
			"detail": "Token inválido o expirado",
		}}})
		return
	}

	if time.Since(pr.CreatedAt) > 2*time.Hour {
		_ = providers.DeletePasswordReset(req.Token)
		context.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{
			"status": "400",
			"title":  "Token Expired",
			"detail": "Token expirado",
		}}})
		return
	}

	user, err := providers.GetUserByEmail(pr.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{
			"status": "400",
			"title":  "User Not Found",
			"detail": "Usuario no encontrado",
		}}})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error al encriptar contraseña",
		}}})
		return
	}

	update := models.UserStruct{ID: user.ID, Password: string(hashedPassword)}
	err = providers.UpdateUser(update)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "No se pudo actualizar la contraseña",
		}}})
		return
	}

	_ = providers.DeletePasswordReset(req.Token)
	context.JSON(http.StatusOK, gin.H{"message": "Contraseña restablecida"})
}
