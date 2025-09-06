package requests

import "semita/core/validators"

type ForgotPasswordRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Email string `json:"email"`
		} `json:"attributes"`
	} `json:"data"`
}

func (request *ForgotPasswordRequest) Rules() *validators.Validator {
	validator := validators.New()

	validator.Field("data").Required()
	validator.Field("data.type").Required().String().In("users")
	validator.Field("data.attributes").Required()
	validator.Field("data.attributes.email").Required().Email().Exists("users", "email")

	return validator
}

func (request *ForgotPasswordRequest) Messages() map[string]string {
	return map[string]string{
		"data.required":                  "Los datos son obligatorios",
		"data.type.required":             "El tipo es obligatorio",
		"data.type.in":                   "El tipo debe ser 'users'",
		"data.attributes.required":       "Los atributos son obligatorios",
		"data.attributes.email.required": "El email es obligatorio",
		"data.attributes.email.email":    "Debe ser un email válido",
		"data.attributes.email.exists":   "El email no está registrado",
	}
}
