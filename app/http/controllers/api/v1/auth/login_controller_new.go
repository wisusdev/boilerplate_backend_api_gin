package auth

import (
	"net/http"
	"semita/app/data/models"
	"semita/app/http/requests"
	"semita/app/http/resources"
	"semita/core/oauth/oauth_models"
	"semita/core/validators"
	"semita/core/validators/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginNew controlador de login usando el nuevo sistema de validaciones
func LoginNew(context *gin.Context) {
	var request requests.LoginRequestNew

	// Validar usando el nuevo sistema - manejo automático de errores
	if err := validators.Validate(context, &request); err != nil {
		return // Los errores se manejan automáticamente con formato JSON API
	}

	// Buscar usuario por email
	storedUser, err := models.GetUserByEmail(request.Email)
	if err != nil {
		context.JSON(http.StatusUnauthorized, validators.ValidationResponse{
			Errors: []validators.ValidationErrorResponse{{
				Title:  "Unauthorized",
				Detail: "Invalid email or password",
				Meta: validators.ValidationErrorMeta{
					Field:   "email",
					Rule:    "credentials",
					Message: "Email o contraseña incorrectos",
				},
			}},
		})
		return
	}

	// Verificar contraseña
	errPassword := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(request.Password))
	if errPassword != nil {
		context.JSON(http.StatusUnauthorized, validators.ValidationResponse{
			Errors: []validators.ValidationErrorResponse{{
				Title:  "Unauthorized",
				Detail: "Invalid email or password",
				Meta: validators.ValidationErrorMeta{
					Field:   "password",
					Rule:    "credentials",
					Message: "Email o contraseña incorrectos",
				},
			}},
		})
		return
	}

	// Generar token OAuth
	clients, err := oauth_models.GetAllClients()
	if err != nil || len(clients) == 0 {
		context.JSON(http.StatusInternalServerError, validators.ValidationResponse{
			Errors: []validators.ValidationErrorResponse{{
				Title:  "Server Error",
				Detail: "No OAuth client available",
				Meta: validators.ValidationErrorMeta{
					Field:   "server",
					Rule:    "oauth_config",
					Message: "Error de configuración del servidor",
				},
			}},
		})
		return
	}

	client := clients[0]
	token, err := oauth_models.CreateToken(int64(storedUser.ID), client.ID, "")
	if err != nil {
		context.JSON(http.StatusInternalServerError, validators.ValidationResponse{
			Errors: []validators.ValidationErrorResponse{{
				Title:  "Server Error",
				Detail: "Error generating OAuth token",
				Meta: validators.ValidationErrorMeta{
					Field:   "server",
					Rule:    "token_generation",
					Message: "Error al generar token de acceso",
				},
			}},
		})
		return
	}

	// Respuesta exitosa
	resource := resources.NewAuthResource(uint(storedUser.ID), storedUser.FirstName+" "+storedUser.LastName, storedUser.Email, token.AccessToken)
	response := resources.NewAuthLoginResponse(resource, token.RefreshToken, 86400, token.Scopes)
	context.JSON(http.StatusOK, response)
}

// RegisterNew controlador de registro usando el nuevo sistema
func RegisterNew(context *gin.Context) {
	var request requests.RegisterRequestNew

	if err := validators.Validate(context, &request); err != nil {
		return
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, validators.ValidationResponse{
			Errors: []validators.ValidationErrorResponse{{
				Title:  "Server Error",
				Detail: "Error hashing password",
				Meta: validators.ValidationErrorMeta{
					Field:   "password",
					Rule:    "hash",
					Message: "Error al procesar la contraseña",
				},
			}},
		})
		return
	}

	// Crear estructura de usuario
	user := struct {
		Name     string
		Email    string
		Password string
	}{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	// Guardar usuario (simulado)
	// En tu implementación real usarías models.StoreUser(user)

	context.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado exitosamente",
		"data": gin.H{
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// ForgotPasswordNew controlador para solicitar reset de contraseña
func ForgotPasswordNew(context *gin.Context) {
	var request requests.ForgotPasswordRequestNew

	if err := validators.Validate(context, &request); err != nil {
		return
	}

	// Generar token de reset
	// token := utils.GenerateResetToken(request.Email)
	// resetURL := "http://" + utils.GetEnv("APP_URL") + "/auth/reset-password?token=" + token

	// Guardar token en BD
	// _ = models.CreatePasswordReset(request.Email, token)

	// Enviar email
	// errorSendEmail := notifications.SendPasswordReset(request.Email, resetURL)

	context.JSON(http.StatusOK, gin.H{
		"message": "Si el email existe, recibirás un enlace para restablecer tu contraseña",
	})
}

// ResetPasswordNew controlador para restablecer contraseña
func ResetPasswordNew(context *gin.Context) {
	var request requests.ResetPasswordRequestNew

	if err := validators.Validate(context, &request); err != nil {
		return
	}

	// Verificar token válido
	// passwordReset, err := models.GetPasswordResetByToken(request.Token)
	// if err != nil {
	//     context.JSON(http.StatusBadRequest, validators.ValidationResponse{...})
	//     return
	// }

	// Hash nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, validators.ValidationResponse{
			Errors: []validators.ValidationErrorResponse{{
				Title:  "Server Error",
				Detail: "Error hashing password",
				Meta: validators.ValidationErrorMeta{
					Field:   "password",
					Rule:    "hash",
					Message: "Error al procesar la contraseña",
				},
			}},
		})
		return
	}

	// Actualizar contraseña del usuario
	// user, err := models.GetUserByEmail(request.Email)
	// if err != nil { ... }

	// update := structs.UpdateUserStruct{
	//     ID: user.ID,
	//     Password: string(hashedPassword),
	// }
	// err = models.UpdateUser(update)

	// Eliminar token usado
	// _ = models.DeletePasswordReset(request.Token)

	_ = hashedPassword // Usar la variable para evitar error de compilación

	context.JSON(http.StatusOK, gin.H{
		"message": "Contraseña actualizada exitosamente",
	})
}

// Ejemplo de cómo configurar las rutas con ambos sistemas (migración gradual)
func SetupAuthRoutes(router *gin.Engine) {
	api := router.Group("/api/v1/auth")

	// Sistema actual (mantener durante migración)
	api.POST("/login", LoginOld)       // Tu función actual
	api.POST("/register", RegisterOld) // Tu función actual

	// Nuevo sistema (versiones mejoradas)
	api.POST("/login/v2", LoginNew)       // Nueva versión
	api.POST("/register/v2", RegisterNew) // Nueva versión
	api.POST("/forgot-password/v2", ForgotPasswordNew)
	api.POST("/reset-password/v2", ResetPasswordNew)

	// O usando middleware (aún más limpio)
	api.POST("/login/middleware",
		middleware.ValidateJSON(&requests.LoginRequestNew{}),
		func(c *gin.Context) {
			// Request ya validado está en el contexto
			// Implementar lógica aquí...
		})
}

// Funciones placeholder para el ejemplo
func LoginOld(c *gin.Context) {
	// Tu implementación actual
}

func RegisterOld(c *gin.Context) {
	// Tu implementación actual
}
