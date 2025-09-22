package requests

import (
	"boilerplate_backend_api_gin/core/validators"
)

type LoginRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"attributes"`
	} `json:"data"`
}

func (request *LoginRequest) Rules() *validators.Validator {
	validator := validators.New()

	// Validar estructura JSON:API
	validator.Field("data").Required()
	validator.Field("data.type").Required().String().In("users")
	validator.Field("data.attributes").Required()
	validator.Field("data.attributes.email").Required().Email().Exists("users", "email")
	validator.Field("data.attributes.password").Required().Min(6)

	return validator
}

func (request *LoginRequest) Messages() map[string]string {
	return map[string]string{
		"data.required":                     "Los datos son obligatorios",
		"data.type.required":                "El tipo es obligatorio",
		"data.type.in":                      "El tipo debe ser 'users'",
		"data.attributes.required":          "Los atributos son obligatorios",
		"data.attributes.email.required":    "El email es obligatorio",
		"data.attributes.email.email":       "Debe ser un email válido",
		"data.attributes.email.exists":      "El email no está registrado",
		"data.attributes.password.required": "La contraseña es obligatoria",
		"data.attributes.password.min":      "La contraseña debe tener al menos 6 caracteres",
	}
}
