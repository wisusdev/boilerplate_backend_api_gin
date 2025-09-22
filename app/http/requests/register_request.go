package requests

import "boilerplate_backend_api_gin/core/validators"

type RegisterRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Username        string `json:"username"`
			Firstname       string `json:"first_name"`
			Lastname        string `json:"last_name"`
			Email           string `json:"email"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"password_confirmation"`
		} `json:"attributes"`
	} `json:"data"`
}

func (request *RegisterRequest) Rules() *validators.Validator {
	registerValidator := validators.New()

	registerValidator.Field("data").Required()
	registerValidator.Field("data.type").Required().String().In("users")
	registerValidator.Field("data.attributes").Required()
	registerValidator.Field("data.attributes.username").Required().Min(2).Max(50).AlphaNum().Unique("users", "username")
	registerValidator.Field("data.attributes.first_name").Required().Min(2).Max(50).Alpha()
	registerValidator.Field("data.attributes.last_name").Required().Min(2).Max(50).Alpha()
	registerValidator.Field("data.attributes.email").Required().Email().Unique("users", "email")
	registerValidator.Field("data.attributes.password").Required().Min(6)
	registerValidator.Field("data.attributes.password_confirmation").Required().Min(6).Same("data.attributes.password")

	return registerValidator
}

func (request *RegisterRequest) Messages() map[string]string {
	return map[string]string{
		"data.required":                                  "Los datos son obligatorios",
		"data.type.required":                             "El tipo es obligatorio",
		"data.type.in":                                   "El tipo debe ser 'users'",
		"data.attributes.required":                       "Los atributos son obligatorios",
		"data.attributes.username.required":              "El nombre de usuario es obligatorio",
		"data.attributes.username.min":                   "El nombre de usuario debe tener al menos 2 caracteres",
		"data.attributes.username.max":                   "El nombre de usuario debe tener máximo 50 caracteres",
		"data.attributes.username.alpha":                 "El nombre de usuario solo puede contener letras",
		"data.attributes.first_name.required":            "El nombre es obligatorio",
		"data.attributes.first_name.min":                 "El nombre debe tener al menos 2 caracteres",
		"data.attributes.first_name.max":                 "El nombre debe tener máximo 50 caracteres",
		"data.attributes.first_name.alpha":               "El nombre solo puede contener letras",
		"data.attributes.last_name.required":             "El apellido es obligatorio",
		"data.attributes.last_name.min":                  "El apellido debe tener al menos 2 caracteres",
		"data.attributes.last_name.max":                  "El apellido debe tener máximo 50 caracteres",
		"data.attributes.last_name.alpha":                "El apellido solo puede contener letras",
		"data.attributes.email.required":                 "El email es obligatorio",
		"data.attributes.email.email":                    "Debe ser un email válido",
		"data.attributes.email.unique":                   "El email ya está registrado",
		"data.attributes.password.required":              "La contraseña es obligatoria",
		"data.attributes.password.min":                   "La contraseña debe tener al menos 6 caracteres",
		"data.attributes.password_confirmation.required": "La confirmación de contraseña es obligatoria",
		"data.attributes.password_confirmation.min":      "La confirmación de contraseña debe tener al menos 6 caracteres",
		"data.attributes.password_confirmation.same":     "La confirmación de contraseña debe coincidir con la contraseña",
	}
}
