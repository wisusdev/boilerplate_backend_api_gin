package requests

import "boilerplate_backend_api_gin/core/validators"

type ResetPasswordRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Token                string `json:"token"`
			Password             string `json:"password"`
			PasswordConfirmation string `json:"password_confirmation"`
		} `json:"attributes"`
	} `json:"data"`
}

func (r *ResetPasswordRequest) Rules() *validators.Validator {
	validator := validators.New()

	validator.Field("data").Required()
	validator.Field("data.type").Required().String().In("reset-password")
	validator.Field("data.attributes").Required()
	validator.Field("data.attributes.token").Required()
	validator.Field("data.attributes.password").Required().Min(6)
	validator.Field("data.attributes.password_confirmation").Required().Min(6).Same("data.attributes.password")

	return validator
}

func (r *ResetPasswordRequest) Messages() map[string]string {
	return map[string]string{
		"data.required":                                  "Los datos son obligatorios",
		"data.type.required":                             "El tipo es obligatorio",
		"data.type.in":                                   "El tipo debe ser 'reset-password'",
		"data.attributes.required":                       "Los atributos son obligatorios",
		"data.attributes.token.required":                 "El token es obligatorio",
		"data.attributes.password.required":              "La contraseña es obligatoria",
		"data.attributes.password.min":                   "La contraseña debe tener al menos 6 caracteres",
		"data.attributes.password.confirmed":             "Las contraseñas no coinciden",
		"data.attributes.password_confirmation.required": "La confirmación de la contraseña es obligatoria",
		"data.attributes.password_confirmation.min":      "La confirmación de la contraseña debe tener al menos 6 caracteres",
		"data.attributes.password_confirmation.same":     "La confirmación de la contraseña debe coincidir con la contraseña",
	}
}
