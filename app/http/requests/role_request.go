package requests

import (
	"boilerplate_backend_api_gin/core/validators"
)

type RoleRequest struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id,omitempty"`
		Attributes struct {
			Name        string   `json:"name"`
			Permissions []string `json:"permissions,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

func (request *RoleRequest) Rules() *validators.Validator {
	validator := validators.New()

	// Validar estructura JSON:API
	validator.Field("data").Required()
	validator.Field("data.type").Required().String().In("roles")
	validator.Field("data.attributes").Required()
	validator.Field("data.attributes.name").Required().String().Min(2).Max(100)

	// Validar unicidad del nombre
	// Para actualización, la lógica de unicidad se manejará en el controlador
	validator.Field("data.attributes.name").Unique("roles", "name")

	// Validar que todos los permisos existan en la base de datos
	validator.Field("data.attributes.permissions").ArrayExists("permissions", "name")

	return validator
}

func (request *RoleRequest) Messages() map[string]string {
	return map[string]string{
		"data.required":                            "Los datos son obligatorios",
		"data.type.required":                       "El tipo es obligatorio",
		"data.type.in":                             "El tipo debe ser 'roles'",
		"data.attributes.required":                 "Los atributos son obligatorios",
		"data.attributes.name.required":            "El nombre del rol es obligatorio",
		"data.attributes.name.string":              "El nombre debe ser una cadena de texto",
		"data.attributes.name.min":                 "El nombre debe tener al menos 2 caracteres",
		"data.attributes.name.max":                 "El nombre no puede tener más de 100 caracteres",
		"data.attributes.name.unique":              "Ya existe un rol con este nombre",
		"data.attributes.permissions.array_exists": "Algunos permisos no existen en el sistema",
	}
}

// IsUpdate indica si es una operación de actualización
func (request *RoleRequest) IsUpdate() bool {
	return request.Data.ID != ""
}
