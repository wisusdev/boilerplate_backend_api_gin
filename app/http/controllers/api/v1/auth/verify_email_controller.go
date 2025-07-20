package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"semita/app/data/models"
	"semita/app/notifications"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ResendEmailVerify(context *gin.Context) {
	// Simulación: obtener usuario autenticado (en real, usar JWT o sesión)
	userId := 1 // TODO: obtener del contexto real
	user, err := models.GetUserByID(strconv.Itoa(userId))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	// Generar hash de verificación
	hash := generateEmailVerificationHash(user.ID, user.Email)
	verifyURL := "/auth/email/verify/" + strconv.Itoa(user.ID) + "/" + hash
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
	user, err := models.GetUserByID(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	// Validar hash
	if hash != generateEmailVerificationHash(user.ID, user.Email) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Hash inválido"})
		return
	}
	// Marcar email como verificado
	err = models.MarkEmailVerified(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo verificar el email"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Email verificado correctamente"})
}

func generateEmailVerificationHash(userID int, email string) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d:%s", userID, email)))
	return hex.EncodeToString(h.Sum(nil))
}
