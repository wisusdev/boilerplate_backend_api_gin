package requests

import (
	"boilerplate_backend_api_gin/core/validators"
)

type UserRequest struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id,omitempty"`
		Attributes struct {
			Username             string   `json:"username"`
			FirstName            string   `json:"first_name"`
			LastName             string   `json:"last_name"`
			Email                string   `json:"email"`
			Password             string   `json:"password,omitempty"`
			PasswordConfirmation string   `json:"password_confirmation,omitempty"`
			Roles                []string `json:"roles,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

func (request *UserRequest) Rules() *validators.Validator {
	validator := validators.New()

	// Validar estructura JSON:API
	validator.Field("data").Required()
	validator.Field("data.type").Required().String().In("users")
	validator.Field("data.attributes").Required()

	// Validar campos requeridos
	validator.Field("data.attributes.username").Required().String().Min(3).Max(50)
	validator.Field("data.attributes.first_name").Required().String().Min(2).Max(100)
	validator.Field("data.attributes.last_name").Required().String().Min(2).Max(100)
	validator.Field("data.attributes.email").Required().Email()

	// Para creación, password es requerido
	if !request.IsUpdate() {
		validator.Field("data.attributes.password").Required().String().Min(6)
		validator.Field("data.attributes.password_confirmation").Required().String().Confirmed()
	} else {
		// Para actualización, si hay password debe cumplir las reglas
		if request.Data.Attributes.Password != "" {
			validator.Field("data.attributes.password").String().Min(6)
			validator.Field("data.attributes.password_confirmation").String().Confirmed()
		}
	}

	// Validar unicidad de campos
	validator.Field("data.attributes.username").Unique("users", "username")
	validator.Field("data.attributes.email").Unique("users", "email")

	// Validar que los roles existan
	validator.Field("data.attributes.roles").ArrayExists("roles", "name")

	return validator
}

func (request *UserRequest) Messages() map[string]string {
	return map[string]string{
		"data.required":                                   "Los datos son obligatorios",
		"data.type.required":                              "El tipo es obligatorio",
		"data.type.in":                                    "El tipo debe ser 'users'",
		"data.attributes.required":                        "Los atributos son obligatorios",
		"data.attributes.username.required":               "El nombre de usuario es obligatorio",
		"data.attributes.username.string":                 "El nombre de usuario debe ser una cadena de texto",
		"data.attributes.username.min":                    "El nombre de usuario debe tener al menos 3 caracteres",
		"data.attributes.username.max":                    "El nombre de usuario no puede tener más de 50 caracteres",
		"data.attributes.username.unique":                 "Ya existe un usuario con este nombre de usuario",
		"data.attributes.first_name.required":             "El nombre es obligatorio",
		"data.attributes.first_name.string":               "El nombre debe ser una cadena de texto",
		"data.attributes.first_name.min":                  "El nombre debe tener al menos 2 caracteres",
		"data.attributes.first_name.max":                  "El nombre no puede tener más de 100 caracteres",
		"data.attributes.last_name.required":              "El apellido es obligatorio",
		"data.attributes.last_name.string":                "El apellido debe ser una cadena de texto",
		"data.attributes.last_name.min":                   "El apellido debe tener al menos 2 caracteres",
		"data.attributes.last_name.max":                   "El apellido no puede tener más de 100 caracteres",
		"data.attributes.email.required":                  "El email es obligatorio",
		"data.attributes.email.email":                     "Debe ser un email válido",
		"data.attributes.email.unique":                    "Ya existe un usuario con este email",
		"data.attributes.password.required":               "La contraseña es obligatoria",
		"data.attributes.password.string":                 "La contraseña debe ser una cadena de texto",
		"data.attributes.password.min":                    "La contraseña debe tener al menos 6 caracteres",
		"data.attributes.password_confirmation.required":  "La confirmación de contraseña es obligatoria",
		"data.attributes.password_confirmation.string":    "La confirmación de contraseña debe ser una cadena de texto",
		"data.attributes.password_confirmation.confirmed": "La confirmación de contraseña no coincide",
		"data.attributes.roles.array_exists":              "Algunos roles no existen en el sistema",
	}
}

// IsUpdate indica si es una operación de actualización
func (request *UserRequest) IsUpdate() bool {
	return request.Data.ID != ""
}
