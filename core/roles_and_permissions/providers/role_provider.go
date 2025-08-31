package providers

import (
	"semita/core/roles_and_permissions/models"
	"semita/core/roles_and_permissions/repositories"
)

// GetAllRoles obtiene todos los roles
func GetAllRoles() ([]models.RoleStruct, error) {
	return repositories.GetAllRoles()
}

// GetRoleByID obtiene un rol por su ID
func GetRoleByID(id int) (*models.RoleStruct, error) {
	return repositories.GetRoleByID(id)
}

// GetRoleByName obtiene un rol por su nombre
func GetRoleByName(name string, guardName string) (*models.RoleStruct, error) {
	return repositories.GetRoleByName(name, guardName)
}

// CreateRole crea un nuevo rol
func CreateRole(role models.CreateRoleStruct) (*models.RoleStruct, error) {
	return repositories.CreateRole(role)
}

// UpdateRole actualiza un rol existente
func UpdateRole(id int, role models.CreateRoleStruct) (*models.RoleStruct, error) {
	return repositories.UpdateRole(id, role)
}

// DeleteRole elimina un rol
func DeleteRole(id int) error {
	return repositories.DeleteRole(id)
}

// GetUserRoles obtiene todos los roles de un usuario
func GetUserRoles(userID string) ([]models.RoleStruct, error) {
	return repositories.GetUserRoles(userID)
}

// AssignRoleToUser asigna un rol a un usuario
func AssignRoleToUser(userID string, roleID string) error {
	return repositories.AssignRoleToUser(userID, roleID)
}

// RevokeRoleFromUser revoca un rol de un usuario
func RevokeRoleFromUser(userID string, roleID string) error {
	return repositories.RevokeRoleFromUser(userID, roleID)
}

// UserHasRole verifica si un usuario tiene un rol específico
func UserHasRole(userID string, roleID string) (bool, error) {
	return repositories.UserHasRole(userID, roleID)
}

// UserHasRoleByName verifica si un usuario tiene un rol por nombre
func UserHasRoleByName(userID string, roleName string, guardName string) (bool, error) {
	return repositories.UserHasRoleByName(userID, roleName, guardName)
}

// UserHasAnyRole verifica si un usuario tiene al menos uno de los roles especificados
func UserHasAnyRole(userID string, roleNames []string, guardName string) (bool, error) {
	return repositories.UserHasAnyRole(userID, roleNames, guardName)
}

// UserHasAllRoles verifica si un usuario tiene todos los roles especificados
func UserHasAllRoles(userID string, roleNames []string, guardName string) (bool, error) {
	return repositories.UserHasAllRoles(userID, roleNames, guardName)
}
