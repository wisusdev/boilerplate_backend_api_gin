package web

import (
	"fmt"
	"net/http"
	"semita/app/data/models"
	"semita/app/data/structs"
	"semita/app/notifications"
	"semita/config"
	"semita/core/helpers"
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
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Email and password are required")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	user := structs.LoginUserStruct{
		Email:    email,
		Password: password,
	}

	storedUser, err := models.GetUserByEmail(user.Email)
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

	var userData = helpers.UserSessionStruct{
		ID:        storedUser.ID,
		FirstName: storedUser.FirstName,
		LastName:  storedUser.LastName,
		Username:  storedUser.Username,
		Avatar:    helpers.StringToNullString(storedUser.Avatar),
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

func AuthRegisterPost(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	if name == "" || email == "" || password == "" || confirmPassword == "" {
		c.String(http.StatusBadRequest, "Name, Email, password, and confirm password are required")
		return
	}

	if password != confirmPassword {
		c.String(http.StatusBadRequest, "Passwords do not match")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error encrypting password")
		return
	}

	user := structs.StoreUserStruct{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	errorStore := models.StoreUser(user)
	if errorStore != nil {
		helpers.Logs("ERROR", fmt.Sprintf("Error saving user: %v", errorStore))
		c.String(http.StatusInternalServerError, "Error saving user to the database")
		return
	}

	c.Redirect(http.StatusSeeOther, "/auth/login")
	c.Abort()
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
	_ = models.CreatePasswordReset(email, token) // Guardar token en BD
	errorSendEmail := notifications.SendPasswordReset(email, resetURL)

	if errorSendEmail != nil {
		helpers.Logs("ERROR", errorSendEmail.Error())
		fmt.Println("Error sending password reset email:", errorSendEmail)
		return
	}

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
		return
	}

	if password != confirmPassword {
		helpers.CreateFlashNotification(context.Writer, context.Request, "warning", "Passwords do not match")
		return
	}

	passwordResetByToken, err := models.GetPasswordResetByToken(token)
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
		_ = models.DeletePasswordReset(token)
		helpers.Logs("INFO", fmt.Sprintf("Token expirado. Creado hace: %v", timeSince))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "Token expirado. Por favor, solicita un nuevo enlace de restablecimiento.")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	user, err := models.GetUserByEmail(passwordResetByToken.Email)
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

	update := structs.UpdateUserStruct{ID: user.ID, Name: user.FirstName + " " + user.LastName, Email: user.Email, Password: string(hashedPassword)}
	err = models.UpdateUser(update)

	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("No se pudo actualizar la contraseña: %v", err))
		helpers.CreateFlashNotification(context.Writer, context.Request, "error", "No se pudo actualizar la contraseña")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	// Eliminar el token después de usarlo exitosamente
	_ = models.DeletePasswordReset(token)
	helpers.Logs("INFO", "Contraseña restablecida exitosamente")

	helpers.CreateFlashNotification(context.Writer, context.Request, "success", "Contraseña actualizada exitosamente!")
	context.Redirect(http.StatusSeeOther, "/auth/login")
	context.Abort()
}
