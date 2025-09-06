package auth

import (
	"net/http"
	"semita/app/data/providers"
	"semita/app/notifications"
	"semita/config"
	"semita/core/helpers"

	"github.com/gin-gonic/gin"
)

func ResendEmailVerify(context *gin.Context) {
	// Obtener user_id del token OAuth (establecido por AuthMiddleware)
	userIdInterface, exists := context.Get("user_id")
	if !exists {
		helpers.Logs("ERROR", "user_id not found in context")
		context.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}

	userId, ok := userIdInterface.(string)
	if !ok {
		helpers.Logs("ERROR", "user_id is not a string")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	user, err := providers.GetUserByID(userId)
	if err != nil {
		helpers.Logs("ERROR", "Error recuperando usuario: "+err.Error())
		context.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	// Generar hash de verificación
	hash := helpers.GenerateEmailVerificationHash(user.ID, user.Email)
	verifyURL := "http://" + config.AppConfig().Url + "/auth/email/verify/" + user.ID + "/" + hash
	err = notifications.SendEmailVerification(user.Email, verifyURL)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo enviar el correo de verificación"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Correo de verificación enviado", "verify_url": verifyURL})
}

func VerifyEmail(context *gin.Context) {
	id := context.Param("id")
	hash := context.Param("hash")

	if id == "" || hash == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "ID y hash requeridos"})
		return
	}

	// Buscar usuario por ID
	user, err := providers.GetUserByID(id)
	if err != nil {
		helpers.Logs("ERROR", "Error recuperando usuario: "+err.Error())
		context.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if user.EmailVerifiedAt != nil {
		context.JSON(http.StatusOK, gin.H{"message": "Email ya verificado"})
		return
	}

	// Marcar email como verificado
	err = providers.MarkEmailVerified(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo verificar el email"})
		return
	}
	
	context.JSON(http.StatusOK, gin.H{"message": "Email verificado correctamente"})
}
