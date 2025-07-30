package web

import (
	"fmt"
	"net/http"
	"semita/app/data/models"
	"semita/app/data/providers"
	"semita/app/notifications"
	"semita/config"
	"semita/core/helpers"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AuthLogin(context *gin.Context) {
	helpers.View(context, "auth/login.html", "Login", nil)
}

func AuthLoginPost(context *gin.Context) {
	email := context.PostForm("email")
	password := context.PostForm("password")

	if email == "" || password == "" {
		helpers.Logs("ERROR", "Email and password are required")
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "El correo y la contraseña son requeridos")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	user := models.UserStruct{
		Email:    email,
		Password: password,
	}

	storedUser, err := providers.GetUserByEmail(user.Email)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error retrieving user: %v", err))
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Invalid email or password")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	errPassword := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if errPassword != nil {
		helpers.Logs("ERROR", "Invalid password")
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Invalid email or password")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	// Check if email is verified
	if storedUser.EmailVerifiedAt == nil && config.AppConfig().MustVerifyEmail {
		urlResendVerification := "http://" + config.AppConfig().Url + "/auth/resend-verification?email=" + storedUser.Email
		flashMsg := `Por favor verifica tu correo antes de iniciar sesión. <a href="` + urlResendVerification + `">Reenviar correo de verificación</a>`
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", flashMsg)
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	var userData = helpers.UserSessionStruct{
		ID:        storedUser.ID,
		FirstName: storedUser.FirstName,
		LastName:  storedUser.LastName,
		Username:  storedUser.Username,
		Avatar:    storedUser.Avatar,
		Language:  helpers.StringToNullString(storedUser.Language),
		Email:     storedUser.Email,
	}

	sessionLoginError := helpers.LoginUserSession(context.Writer, context.Request, userData)
	if sessionLoginError != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error creating user core_session: %v", sessionLoginError))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Error creating user core_session")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	helpers.CreateFlashNotification(context.Writer, context.Request, "success", "Login successful!")
	context.Redirect(http.StatusSeeOther, "/")
	context.Abort()
}

func AuthLogout(c *gin.Context) {
	sessionLogoutError := helpers.LogoutUserSession(c.Writer, c.Request)
	if sessionLogoutError != nil {
		c.String(http.StatusInternalServerError, "Error logging out")
		return
	}

	helpers.CreateFlashNotification(c.Writer, c.Request, "success", "Logout successful!")
	c.Redirect(http.StatusSeeOther, "/")
	c.Abort()
}

func AuthRegister(context *gin.Context) {
	helpers.View(context, "auth/register.html", "Register", nil)
}

func AuthRegisterPost(context *gin.Context) {
	firstName := context.PostForm("first_name")
	lastName := context.PostForm("last_name")
	username := context.PostForm("username")
	email := context.PostForm("email")
	password := context.PostForm("password")
	confirmPassword := context.PostForm("confirm_password")

	if firstName == "" || lastName == "" || username == "" || email == "" || password == "" || confirmPassword == "" {
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "All fields are required")
		context.Redirect(http.StatusSeeOther, "/auth/register")
		context.Abort()
		return
	}

	if password != confirmPassword {
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Passwords do not match")
		context.Redirect(http.StatusSeeOther, "/auth/register")
		context.Abort()
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error encrypting password: %v", err))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Lo siento, hubo un error al encriptar la contraseña")
		context.Redirect(http.StatusSeeOther, "/auth/register")
		context.Abort()
		return
	}

	user := models.UserStruct{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
	}

	userStore, errorStore := providers.StoreUser(user)
	if errorStore != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error saving user: %v", errorStore))
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Lo siento, hubo un error al guardar el usuario")
		context.Redirect(http.StatusSeeOther, "/auth/register")
		context.Abort()
		return
	}

	if config.AppConfig().MustVerifyEmail {
		verifyURL := "http://" + config.AppConfig().Url + "/auth/email/verify/" + strconv.Itoa(userStore.ID) + "/" + helpers.GenerateEmailVerificationHash(userStore.ID, userStore.Email)
		err = notifications.SendEmailVerification(userStore.Email, verifyURL)

		if err != nil {
			helpers.Logs("ERROR", fmt.Sprintf("Error sending verification email: %v", err))
			helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Error sending verification email. Please try again later.")
			context.Redirect(http.StatusSeeOther, "/auth/register")
			context.Abort()
			return
		}

		helpers.CreateFlashNotification(context.Writer, context.Request, "info", "Verification email sent. Please check your inbox.")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	helpers.CreateFlashNotification(context.Writer, context.Request, "success", "User registered successfully!")
	context.Redirect(http.StatusSeeOther, "/auth/login")
	context.Abort()
}

func AuthForgotPassword(context *gin.Context) {
	helpers.View(context, "auth/forgot_password.html", "Recuperar Contraseña", nil)
}

func AuthForgotPasswordPost(context *gin.Context) {
	email := context.PostForm("email")
	if email == "" {
		context.String(http.StatusBadRequest, "Email is required")
		return
	}

	token := helpers.GenerateResetToken(email)
	resetURL := "http://" + config.AppConfig().Url + "/auth/reset-password?token=" + token
	_ = providers.CreatePasswordReset(email, token) // Guardar token en BD
	errorSendEmail := notifications.SendPasswordReset(email, resetURL)

	if errorSendEmail != nil {
		helpers.Logs("ERROR", errorSendEmail.Error())
		fmt.Println("Error sending password reset email:", errorSendEmail)
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Error sending password reset email. Please try again later.")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	helpers.CreateFlashNotification(context.Writer, context.Request, "success", "Password reset link sent successfully")
	context.Redirect(http.StatusSeeOther, "/auth/login")
	context.Abort()
}

func AuthResetPassword(context *gin.Context) {
	var data = map[string]string{
		"token": context.Query("token"),
	}

	helpers.View(context, "auth/reset_password.html", "Restablecer Contraseña", data)
}

func AuthResetPasswordPost(context *gin.Context) {
	token := context.PostForm("token")
	password := context.PostForm("password")
	confirmPassword := context.PostForm("confirm_password")

	if token == "" || password == "" || confirmPassword == "" {
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Token, password, and confirm password are required")
		context.Redirect(http.StatusSeeOther, "/auth/reset-password?token="+token)
		context.Abort()
		return
	}

	if password != confirmPassword {
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Passwords do not match")
		context.Redirect(http.StatusSeeOther, "/auth/reset-password?token="+token)
		context.Abort()
		return
	}

	passwordResetByToken, err := providers.GetPasswordResetByToken(token)
	if err != nil {
		helpers.Logs("ERROR", err.Error())
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Token inválido o expirado")
		context.Redirect(http.StatusSeeOther, "/auth/reset-password?token="+token)
		context.Abort()
		return
	}

	// Usar hora local para ambos tiempos
	now := time.Now()
	tokenCreatedAt := passwordResetByToken.CreatedAt
	timeSince := now.Sub(tokenCreatedAt)

	// Verificar expiración de 2 horas
	if timeSince > 2*time.Hour {
		_ = providers.DeletePasswordReset(token)
		helpers.Logs("INFO", fmt.Sprintf("Token expirado. Creado hace: %v", timeSince))
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Token expirado. Por favor, solicita un nuevo enlace de restablecimiento.")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	user, err := providers.GetUserByEmail(passwordResetByToken.Email)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Usuario no encontrado: %v", err))
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Usuario no encontrado")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error al encriptar contraseña: %v", err))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Error al encriptar contraseña")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	update := models.UserStruct{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Password:  string(hashedPassword),
	}
	err = providers.UpdateUser(update)

	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("No se pudo actualizar la contraseña: %v", err))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "No se pudo actualizar la contraseña")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	// Eliminar el token después de usarlo exitosamente
	_ = providers.DeletePasswordReset(token)
	helpers.Logs("INFO", "Contraseña restablecida exitosamente")

	helpers.CreateFlashNotification(context.Writer, context.Request, "success", "Contraseña actualizada exitosamente!")
	context.Redirect(http.StatusSeeOther, "/auth/login")
	context.Abort()
}

func AuthVerifyEmail(context *gin.Context) {
	id := context.Param("id")
	hash := context.Param("hash")

	if id == "" || hash == "" {
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Invalid verification link")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	user, err := providers.GetUserByID(id)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("User not found: %v", err))
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "User not found")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	if user.EmailVerifiedAt != nil {
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "El correo ya ha sido verificado")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	err = providers.MarkEmailVerified(user.ID)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error marking email as verified: %v", err))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Error marking email as verified")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	helpers.CreateFlashNotification(context.Writer, context.Request, "success", "Email verified successfully!")
	context.Redirect(http.StatusSeeOther, "/auth/login")
	context.Abort()
}

func AuthResendVerification(context *gin.Context) {
	email := context.Query("email")

	if email == "" {
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Correo no válido")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	user, err := providers.GetUserByEmail(email)
	if err != nil || user.EmailVerifiedAt != nil {
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Usuario no encontrado o ya verificado")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	verifyURL := "http://" + config.AppConfig().Url + "/auth/email/verify/" + strconv.Itoa(user.ID) + "/" + helpers.GenerateEmailVerificationHash(user.ID, user.Email)
	err = notifications.SendEmailVerification(user.Email, verifyURL)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error sending verification email: %v", err))
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Error enviando el correo. Intenta más tarde.")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	helpers.CreateFlashNotification(context.Writer, context.Request, "success", "Correo de verificación reenviado. Revisa tu bandeja de entrada.")
	context.Redirect(http.StatusSeeOther, "/auth/login")
	context.Abort()
}
