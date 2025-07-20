package requests

import (
	"semita/core/validators"
)

// LoginRequest valida los datos de login - MIGRADO AL NUEVO SISTEMA
type LoginRequestNew struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequestNew) Rules() *validators.Validator {
	validator := validators.New()
	validator.Field("email").Required().Email()
	validator.Field("password").Required().Min(6)
	return validator
}

func (r *LoginRequestNew) Messages() map[string]string {
	return map[string]string{
		"email.required":    "El email es obligatorio",
		"email.email":       "Debe ser un email válido",
		"password.required": "La contraseña es obligatoria",
		"password.min":      "La contraseña debe tener al menos 6 caracteres",
	}
}

// RegisterRequest valida los datos de registro - MIGRADO AL NUEVO SISTEMA
type RegisterRequestNew struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (r *RegisterRequestNew) Rules() *validators.Validator {
	validator := validators.New()
	validator.Field("name").Required().Min(2).Max(50).Alpha()
	validator.Field("email").Required().Email().Unique("users", "email")
	validator.Field("password").Required().Min(6).Confirmed()
	return validator
}

func (r *RegisterRequestNew) Messages() map[string]string {
	return map[string]string{
		"name.required":      "El nombre es obligatorio",
		"name.min":           "El nombre debe tener al menos 2 caracteres",
		"name.max":           "El nombre no puede tener más de 50 caracteres",
		"name.alpha":         "El nombre solo puede contener letras",
		"email.required":     "El email es obligatorio",
		"email.email":        "Debe ser un email válido",
		"email.unique":       "Este email ya está registrado",
		"password.required":  "La contraseña es obligatoria",
		"password.min":       "La contraseña debe tener al menos 6 caracteres",
		"password.confirmed": "Las contraseñas no coinciden",
	}
}

// ForgotPasswordRequest valida el email para recuperar contraseña - MIGRADO
type ForgotPasswordRequestNew struct {
	Email string `json:"email"`
}

func (r *ForgotPasswordRequestNew) Rules() *validators.Validator {
	validator := validators.New()
	validator.Field("email").Required().Email().Exists("users", "email")
	return validator
}

func (r *ForgotPasswordRequestNew) Messages() map[string]string {
	return map[string]string{
		"email.required": "El email es obligatorio",
		"email.email":    "Debe ser un email válido",
		"email.exists":   "No encontramos una cuenta con este email",
	}
}

// ResetPasswordRequest valida el reseteo de contraseña - MIGRADO
type ResetPasswordRequestNew struct {
	Token                string `json:"token"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (r *ResetPasswordRequestNew) Rules() *validators.Validator {
	validator := validators.New()
	validator.Field("token").Required()
	validator.Field("email").Required().Email().Exists("users", "email")
	validator.Field("password").Required().Min(6).Confirmed()
	return validator
}

func (r *ResetPasswordRequestNew) Messages() map[string]string {
	return map[string]string{
		"token.required":     "El token es obligatorio",
		"email.required":     "El email es obligatorio",
		"email.email":        "Debe ser un email válido",
		"email.exists":       "No encontramos una cuenta con este email",
		"password.required":  "La contraseña es obligatoria",
		"password.min":       "La contraseña debe tener al menos 6 caracteres",
		"password.confirmed": "Las contraseñas no coinciden",
	}
}
